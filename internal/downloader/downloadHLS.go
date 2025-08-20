package downloader

import (
	"bufio"
	"bytes"
	"film-downloader/internal/config"
	"film-downloader/internal/models"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func DownloadHLS(movie models.Movie, cfg *config.Config) error {
	os.MkdirAll(movie.Name, 0755)

	masterBody := downloadFile(movie.Source, cfg.AccessToken)
	defer masterBody.Close()

	masterContent, err := io.ReadAll(masterBody)

	if err != nil {
		return fmt.Errorf("failed to read master playlist: %v", err)
	}

	baseURL, err := url.Parse(movie.Source)

	if err != nil {
		return fmt.Errorf("failed to parse base URL: %v", err)
	}

	var videoM3U8 string
	audioM3U8s := make(map[string]string)
	subtitleM3U8s := make(map[string]string)

	scanner := bufio.NewScanner(bytes.NewReader(masterContent))

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "#EXT-X-STREAM-INF") {
			scanner.Scan()
			videoM3U8 = resolveURL(baseURL, scanner.Text())
		}

		if strings.HasPrefix(line, "#EXT-X-MEDIA") && strings.Contains(line, "TYPE=AUDIO") {
			lang := extractAttr(line, "LANGUAGE")
			uri := extractAttr(line, "URI")

			if uri != "" && lang != "" {
				audioM3U8s[lang] = resolveURL(baseURL, uri)
			}
		}

		if strings.HasPrefix(line, "#EXT-X-MEDIA") && strings.Contains(line, "TYPE=SUBTITLES") {
			lang := extractAttr(line, "LANGUAGE")
			uri := extractAttr(line, "URI")
			if uri != "" && lang != "" {
				subtitleM3U8s[lang] = resolveURL(baseURL, uri)
			}
		}
	}

	if videoM3U8 == "" {
		return fmt.Errorf("no video playlist found")
	}

	downloadMediaPlaylist(movie.Name, videoM3U8, "video/1080p", cfg.AccessToken)

	for lang, audioURL := range audioM3U8s {
		downloadMediaPlaylistWithLang(movie.Name, audioURL, "audio", lang, cfg.AccessToken)
	}

	for lang, subtitleURL := range subtitleM3U8s {
		downloadMediaPlaylistWithLang(movie.Name, subtitleURL, "sub", lang, cfg.AccessToken)
	}

	generateLocalMasterPlaylist(movie.Name, audioM3U8s, subtitleM3U8s)
	return nil
}

func downloadMediaPlaylist(name, playlistURL, folder, accessToken string) []string {
	u, _ := url.Parse(playlistURL)
	resp := downloadFile(playlistURL, accessToken)
	defer resp.Close()
	dir := filepath.Join(name, folder)
	os.MkdirAll(dir, 0755)
	localM3U8 := filepath.Join(dir, "1080p_playlist.m3u8")
	out, err := os.Create(localM3U8)

	if err != nil {
		log.Fatal(err)
	}

	defer out.Close()
	scanner := bufio.NewScanner(resp)
	segmentIndex := 0
	var downloaded []string

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "#") {
			out.WriteString(line + "\n")
			continue
		}

		segmentIndex++
		segmentURL := resolveURL(u, line)
		var fileExt string

		if strings.HasPrefix(folder, "subtitle_") {

			if strings.Contains(line, ".") {
				parts := strings.Split(line, ".")
				fileExt = "." + parts[len(parts)-1]
			} else {
				fileExt = ".vtt"
			}
		} else {
			fileExt = ".ts"
		}

		localSegment := fmt.Sprintf("segment_%03d%s", segmentIndex, fileExt)
		localPath := filepath.Join(dir, localSegment)
		fmt.Printf("Downloading %s...\n", localPath)
		err := saveFile(segmentURL, accessToken, localPath)

		if err != nil {
			log.Fatalf("Failed to download segment: %v", err)
		}

		out.WriteString(localSegment + "\n")
		downloaded = append(downloaded, localSegment)
	}

	return downloaded
}

func downloadMediaPlaylistWithLang(name, playlistURL, folder, lang, accessToken string) []string {
	u, _ := url.Parse(playlistURL)
	resp := downloadFile(playlistURL, accessToken)
	defer resp.Close()

	dir := filepath.Join(name, folder)
	os.MkdirAll(dir, 0755)

	localM3U8 := filepath.Join(dir, fmt.Sprintf("%s_playlist.m3u8", lang))
	out, err := os.Create(localM3U8)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	scanner := bufio.NewScanner(resp)
	segmentIndex := 0
	var downloaded []string

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "#") {
			out.WriteString(line + "\n")
			continue
		}

		segmentIndex++
		segmentURL := resolveURL(u, line)

		var fileExt string
		if strings.HasPrefix(folder, "sub") {
			if strings.Contains(line, ".") {
				parts := strings.Split(line, ".")
				fileExt = "." + parts[len(parts)-1]
			} else {
				fileExt = ".vtt"
			}
		} else {
			fileExt = ".ts"
		}

		localSegment := fmt.Sprintf("%s_segment_%03d%s", lang, segmentIndex, fileExt)
		localPath := filepath.Join(dir, localSegment)
		fmt.Printf("Downloading %s...\n", localPath)
		err := saveFile(segmentURL, accessToken, localPath)

		if err != nil {
			log.Fatalf("Failed to download segment: %v", err)
		}

		out.WriteString(localSegment + "\n")
		downloaded = append(downloaded, localSegment)
	}

	return downloaded
}

func downloadFile(fileURL, auth string) io.ReadCloser {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fileURL, nil)

	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", auth)
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 200 {
		log.Fatalf("Failed to fetch %s: %s", fileURL, resp.Status)
	}
	return resp.Body
}

func saveFile(fileURL, auth, path string) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fileURL, nil)

	if err != nil {
		return err
	}

	req.Header.Set("Authorization", auth)
	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("bad response: %s", resp.Status)
	}

	out, err := os.Create(path)

	if err != nil {
		return err
	}

	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

func resolveURL(base *url.URL, ref string) string {
	url, err := base.Parse(ref)

	if err != nil {
		log.Fatal(err)
	}
	return url.String()
}

func extractAttr(line, key string) string {
	re := regexp.MustCompile(key + `="([^"]+)"`)
	matches := re.FindStringSubmatch(line)

	if len(matches) >= 2 {
		return matches[1]
	}
	return ""
}

func generateLocalMasterPlaylist(movieName string, audioM3U8s, subtitleM3U8s map[string]string) {
	masterPath := filepath.Join(movieName, "master.m3u8")
	out, err := os.Create(masterPath)
	if err != nil {
		log.Printf("Failed to create local master playlist: %v", err)
		return
	}
	defer out.Close()

	out.WriteString("#EXTM3U\n")

	if len(audioM3U8s) > 0 {
		isFirst := true
		for lang := range audioM3U8s {

			defaultStr := "NO"
			if isFirst {
				defaultStr = "YES"
				isFirst = false
			}
			out.WriteString(fmt.Sprintf("#EXT-X-MEDIA:TYPE=AUDIO,GROUP-ID=\"audio\",LANGUAGE=\"%s\",NAME=\"%s\",DEFAULT=%s,AUTOSELECT=YES,URI=\"audio/%s_playlist.m3u8\"\n",
				lang, lang, defaultStr, lang))
		}
		out.WriteString("\n")
	}

	if len(subtitleM3U8s) > 0 {
		for lang := range subtitleM3U8s {
			out.WriteString(fmt.Sprintf("#EXT-X-MEDIA:TYPE=SUBTITLES,GROUP-ID=\"subs\",LANGUAGE=\"%s\",NAME=\"%s\",DEFAULT=NO,URI=\"sub/%s_playlist.m3u8\"\n",
				lang, lang, lang))
		}
		out.WriteString("\n")
	}

	audioGroup := ""
	subsGroup := ""
	if len(audioM3U8s) > 0 {
		audioGroup = ",AUDIO=\"audio\""
	}
	if len(subtitleM3U8s) > 0 {
		subsGroup = ",SUBTITLES=\"subs\""
	}

	out.WriteString(fmt.Sprintf("#EXT-X-STREAM-INF:BANDWIDTH=5128000,RESOLUTION=1920x1080,CODECS=\"avc1.640028,mp4a.40.2\"%s%s\n", audioGroup, subsGroup))
	out.WriteString("video/1080p/1080p_playlist.m3u8\n")

	fmt.Printf("Generated local master playlist: %s\n", masterPath)
}
