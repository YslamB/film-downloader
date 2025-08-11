package main

import (
	"context"
	"film-downloader/internal/config"
	"film-downloader/internal/cron"
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

	cron.CheckDaily(ctx, &wg)
	quit := make(chan os.Signal, 1)
	defer close(quit)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\nReceived shutdown signal...")
	cancel() // Cancel the context to signal goroutines to stop
	wg.Wait()
	log.Println("Shutting down server...")

}
