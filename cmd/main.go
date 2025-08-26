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
	cfg := config.Init()
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	repo := repositories.NewMovieRepository(cfg)
	err := cron.RefreshToken(ctx, cfg, repo)

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		err := cron.GetLastMovies(ctx, cfg, repo, &wg)

		if err != nil {
			log.Fatal(err)
		}
	}()

	// go func() {
	// 	err := cron.CheckDaily(ctx, &wg)

	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }()

	quit := make(chan os.Signal, 1)
	defer close(quit)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Received shutdown signal...")
	wg.Wait()
	cancel() // Cancel the context to signal goroutines to stop
	log.Println("Shutting down server...")

}
