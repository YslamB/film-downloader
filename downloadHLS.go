package main

import (
	"bufio"
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

func DownloadHLS() {
	os.MkdirAll(outputDir, 0755)

	// Step 1: Download the master playlist
	masterBody := downloadFile(baseM3U8URL, authHeader)
	defer masterBody.Close()

	// Step 2: Parse the master playlist and find video and audio playlists
	baseURL, _ := url.Parse(baseM3U8URL)
	var videoM3U8 string
	audioM3U8s := make(map[string]string) // map[lang]url

	scanner := bufio.NewScanner(masterBody)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "#EXT-X-STREAM-INF") {
			// Next line will be the video playlist
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
	downloadMediaPlaylist(videoM3U8, "video")

	for lang, audioURL := range audioM3U8s {
		fmt.Printf("Downloading audio segments (%s)...\n", lang)
		downloadMediaPlaylist(audioURL, "audio_"+lang)
	}

	fmt.Printf("✅ Download complete! You can now use ffmpeg to convert the video:\n")
	fmt.Printf(`ffmpeg -allowed_extensions ALL -i %s/video_local.m3u8 -c copy output.mp4`+"\n", outputDir)
}

// Download and parse media playlist (.m3u8), download all .ts segments
func downloadMediaPlaylist(playlistURL, folder string) []string {
	u, _ := url.Parse(playlistURL)
	resp := downloadFile(playlistURL, authHeader)
	defer resp.Close()

	dir := filepath.Join(outputDir, folder)
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
		err := saveFile(segmentURL, authHeader, localPath)
		if err != nil {
			log.Fatalf("Failed to download segment: %v", err)
		}

		out.WriteString(localSegment + "\n")
		downloaded = append(downloaded, localSegment)
	}

	return downloaded
}

// Download a file with Authorization header
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

// Save a file locally
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

// Resolve relative URI to full URL
func resolveURL(base *url.URL, ref string) string {
	u, err := base.Parse(ref)
	if err != nil {
		log.Fatal(err)
	}
	return u.String()
}

// Extract attribute from EXT-X-MEDIA line
func extractAttr(line, key string) string {
	re := regexp.MustCompile(key + `="([^"]+)"`)
	matches := re.FindStringSubmatch(line)
	if len(matches) >= 2 {
		return matches[1]
	}
	return ""
}
