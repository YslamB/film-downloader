package main

import (
	"context"
	"film-downloader/internal/config"
	"film-downloader/internal/cron"
	"film-downloader/internal/repositories"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	wg := sync.WaitGroup{}
	cfg := config.Init()
	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	repo := repositories.NewMovieRepository(cfg.AccessToken)
	err := cron.DownloadWithID(ctx, "", "", "441340", cfg, repo)
	// err := utils.UploadFolderToMinio(
	// 	"test/fd29222081571ec9eb1df30ec3956cf7/video/480p", "fd29222081571ec9eb1df30ec3956cf7/video/480p",
	// 	cfg.MINIO_BUCKET, cfg.MINIO_ENDPOINT, cfg.MINIO_ACCESS_KEY,
	// 	cfg.MINIO_SECRET_KEY, cfg.MINIO_SECURE, 10,
	// )

	if err != nil {
		log.Fatal(err)
	}

	wg.Done()
	// cron.CheckDaily(ctx, &wg)
	quit := make(chan os.Signal, 1)
	defer close(quit)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Received shutdown signal...")
	cancel() // Cancel the context to signal goroutines to stop
	wg.Wait()
	log.Println("Shutting down server...")

}
