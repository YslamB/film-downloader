package downloader

import (
	"film-downloader/internal/config"
	"film-downloader/internal/models"
	"film-downloader/internal/utils"
	"log"
	"os"
	"os/exec"
	"time"
)

func DownloadMp4(baseM3U8URL models.Movie, outputDir string, cfg *config.Config) error {
	err := os.MkdirAll(outputDir, 0777)

	if err != nil {
		log.Fatalf("failed to create directory: %e", err)
	}

	cmd := exec.Command(
		"ffmpeg",
		"-headers", "authorization:"+cfg.GetAccessToken(),
		"-i", baseM3U8URL.Sources[0].MasterFile,
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

	err = utils.FFmpegProgressHandler(stderr, 2*time.Second)
	if err != nil {
		log.Fatalf("❌ Error reading stderr: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		log.Fatalf("❌ ffmpeg failed: %v", err)
	}

	log.Println("Download completed:", outputDir)

	return nil
}
