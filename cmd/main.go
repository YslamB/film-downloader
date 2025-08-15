package main

import (
	"context"
	"film-downloader/internal/config"
	"film-downloader/internal/cron"
	"film-downloader/internal/database"
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
	db := database.Init(cfg)
	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	err := cron.DownloadWithID(ctx, "", "", "", cfg, db)

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
