package downloader

import (
	"bufio"
	"film-downloader/internal/config"
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

func DownloadHLS(baseM3U8URL string, cfg config.Config) {
	os.MkdirAll("fileName.mp4", 0755)

	masterBody := downloadFile(baseM3U8URL, cfg.AccessToken)
	defer masterBody.Close()

	baseURL, _ := url.Parse(baseM3U8URL)
	var videoM3U8 string
	audioM3U8s := make(map[string]string)

	scanner := bufio.NewScanner(masterBody)
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
	}

	if videoM3U8 == "" {
		log.Fatal("No video playlist found.")
	}

	fmt.Println("Downloading video segments...")
	downloadMediaPlaylist(videoM3U8, "video", cfg.AccessToken)

	for lang, audioURL := range audioM3U8s {
		fmt.Printf("Downloading audio segments (%s)...\n", lang)
		downloadMediaPlaylist(audioURL, "audio_"+lang, cfg.AccessToken)
	}

	fmt.Printf("✅ Download complete! You can now use ffmpeg to convert the video:\n")
	fmt.Printf(`ffmpeg -allowed_extensions ALL -i %s/video_local.m3u8 -c copy output.mp4`+"\n", "fileName.mp4")
}

func downloadMediaPlaylist(playlistURL, folder, accessToken string) []string {
	u, _ := url.Parse(playlistURL)
	resp := downloadFile(playlistURL, accessToken)
	defer resp.Close()

	dir := filepath.Join("fileName.mp4", folder)
	os.MkdirAll(dir, 0755)

	localM3U8 := filepath.Join(dir, "video_local.m3u8")
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
		localSegment := fmt.Sprintf("segment_%03d.ts", segmentIndex)
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
	u, err := base.Parse(ref)
	if err != nil {
		log.Fatal(err)
	}
	return u.String()
}

func extractAttr(line, key string) string {
	re := regexp.MustCompile(key + `="([^"]+)"`)
	matches := re.FindStringSubmatch(line)
	if len(matches) >= 2 {
		return matches[1]
	}
	return ""
}
