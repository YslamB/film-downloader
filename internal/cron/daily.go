package cron

import (
	"context"
	"film-downloader/internal/config"
	"film-downloader/internal/downloader"
	"film-downloader/internal/models"
	"film-downloader/internal/requests"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
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

func DownloadWithID(ctx context.Context, episodeID, seasonID, filmID string, cfg *config.Config, db *pgxpool.Pool) error {
	var movies []models.Movie
	var err error

	if episodeID == "" && seasonID == "" && filmID != "" {
		source, err := requests.GetFilmSourceURL(ctx, filmID, cfg)

		if err != nil {
			return err
		}

		movies = append(movies, source)
	}

	if seasonID != "" {
		movies, err = requests.GetEpisodesWithSeasonID(seasonID, episodeID, cfg)

		if err != nil {
			return err
		}
	}

	fmt.Println("✅ Received Source files...", movies)

	for i := range movies {
		downloader.DownloadHLS(movies[i], cfg)
	}

	return nil
}
