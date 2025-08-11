package downloader

import (
	"bufio"
	"film-downloader/internal/config"
	"film-downloader/internal/models"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

func DownloadMp4(baseM3U8URL models.Movie, outputDir string, cfg config.Config) {
	err := os.MkdirAll(outputDir, 0777)
	if err != nil {
		log.Fatalf("failed to create directory: %e", err)
	}

	cmd := exec.Command(
		"ffmpeg",
		"-headers", "authorization:"+cfg.AccessToken,
		"-i", baseM3U8URL.Source,
		"-map", "0:v", "-map", "0:a",
		"-c", "copy",
		"-bsf:a", "aac_adtstoasc",
		outputDir+"/"+baseM3U8URL.Name+".mp4",
	)

	stderr, err := cmd.StderrPipe()

	if err != nil {
		log.Fatalf("❌ Error getting stderr pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatalf("❌ Error starting ffmpeg: %v", err)
	}

	scanner := bufio.NewScanner(stderr)

	// Regular expressions to parse info like time, size, speed
	timeRe := regexp.MustCompile(`time=(\d+:\d+:\d+\.\d+)`)
	speedRe := regexp.MustCompile(`speed=([\d\.]+)x`)
	sizeRe := regexp.MustCompile(`size=\s*([\d\.]+)kB`)

	lastLogTime := time.Now()

	for scanner.Scan() {
		line := scanner.Text()

		if time.Since(lastLogTime) >= 2*time.Second {
			timeMatch := timeRe.FindStringSubmatch(line)
			speedMatch := speedRe.FindStringSubmatch(line)
			sizeMatch := sizeRe.FindStringSubmatch(line)

			progress := []string{}

			if len(sizeMatch) > 1 {
				progress = append(progress, "Size: "+sizeMatch[1]+"kB")
			}
			if len(timeMatch) > 1 {
				progress = append(progress, "Time: "+timeMatch[1])
			}
			if len(speedMatch) > 1 {
				progress = append(progress, "Speed: "+speedMatch[1]+"x")
			}

			log.Println(strings.Join(progress, " | "))
			lastLogTime = time.Now()
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("❌ Error reading stderr: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		log.Fatalf("❌ ffmpeg failed: %v", err)
	}

	log.Println("Download completed:", outputDir)
}
