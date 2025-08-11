package cron

import (
	"context"
	"film-downloader/internal/config"
	"film-downloader/internal/downloader"
	"film-downloader/internal/models"
	"film-downloader/internal/requests"
	"fmt"
	"log"
	"sync"
	"time"
)

func CheckDaily(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Stopping daily worker.")
			return
		default:
			// send req to search api, for last 200 items,
			// CheckWithID() in for, get must istall film arrays
			// DownloadWithID(), add wg for each download

			time.Sleep(24 * time.Hour)
		}
	}
}

func CheckWithID() {
	CheckWithStatus()
}

func CheckWithStatus() {

}

func DownloadWithID(episodeID, seasonID, filmID string) {
	var movies []models.Movie
	var err error
	cfg := config.Init()

	if episodeID == "" && seasonID == "" && filmID != "" {
		source, err := requests.GetFilmSourceURL(filmID, cfg)

		if err != nil {
			log.Fatal(err)
		}

		movies = append(movies, source)
	}

	if seasonID != "" {
		movies, err = requests.GetEpisodesWithSeasonID(seasonID, episodeID, cfg)

		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("✅ Received Source files...", movies)

	for i := range movies {
		downloader.DownloadMp4(movies[i], time.Now().Format("2006-01-02"), cfg)
	}

}
