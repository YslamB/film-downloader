package cron

import (
	"context"
	"film-downloader/internal/config"
	"film-downloader/internal/downloader"
	"film-downloader/internal/models"
	"film-downloader/internal/repositories"
	"film-downloader/internal/requests"
	"fmt"
	"os"
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

func DownloadWithID(ctx context.Context, episodeID, seasonID, filmID string, cfg *config.Config, repo *repositories.MovieRepository) error {
	var movies []models.Movie
	var err error

	if episodeID == "" && seasonID == "" && filmID != "" {
		movieSources, err := DownloadMovieSourceWithID(ctx, filmID, cfg, repo)

		if err != nil {
			return err
		}
		movies = append(movies, movieSources...)
	}

	if seasonID != "" {
		fmt.Println("🔍 Checking season with ID:", seasonID)
		movies, err = requests.GetEpisodesSourceWithSeasonID(seasonID, episodeID, cfg)
		time.Sleep(1 * time.Second)

		if err != nil {
			return err
		}
	}

	fmt.Println("✅ Received Source files...", movies)

	for i := range movies {
		err := downloader.DownloadHLS(movies[i], cfg)

		if err != nil {
			return err
		}

		// err = utils.UploadFolderToMinio(
		// 	"temp/"+movies[i].Name, movies[i].Name, cfg.MINIO_BUCKET,
		// 	cfg.MINIO_ENDPOINT, cfg.MINIO_ACCESS_KEY, cfg.MINIO_SECRET_KEY,
		// 	cfg.MINIO_SECURE, cfg.MINIO_WORKERS,
		// )

		// if err != nil {
		// 	return err
		// }

		fileID, err := repo.GetFileID(ctx, movies[i].Name)

		if err != nil {
			return err
		}

		err = repo.CreateMovieFile(ctx, fileID, movies[i].ID)

		if err != nil {
			return err
		}

		fmt.Println("🔍 File ID:", fileID)
		err = os.RemoveAll("temp/" + movies[i].Name)

		if err != nil {
			return err
		}
	}

	return nil
}

func GetLastMovies(ctx context.Context, cfg *config.Config, repo *repositories.MovieRepository) error {
	searchResult, err := requests.GetSearchResults(ctx, 1, cfg)

	if err != nil {
		return fmt.Errorf("failed to get search results from API: %w", err)
	}

	fmt.Printf("📋 Found %d films in search results\n", len(searchResult.Films))

	for i, film := range searchResult.Films {

		fmt.Printf("🎬 Processing film %d/%d: (ID: %d)\n", i+1, len(searchResult.Films), film.ID)
		filmID := fmt.Sprintf("%d", film.ID)

		if film.TypeID == 1 {
			err := DownloadWithID(ctx, "", "", filmID, cfg, repo)

			if err != nil {
				fmt.Printf("❌ Failed to download film %d (ID: %d): %v\n", film.ID, film.ID, err)

				continue
			}
		} else {
			seasons, err := requests.GetSeasonsData(ctx, filmID, cfg)

			if err != nil {
				fmt.Printf("❌ Failed to get seasons data for film %d (ID: %d): %v\n", film.ID, film.ID, err)

				continue
			}

			for _, season := range seasons {
				fmt.Printf("🎬 DownloadWithID season %d/%d: (ID: %d)\n", i+1, len(seasons), season.ID)
				// err := DownloadWithID(ctx, "", fmt.Sprintf("%d", season.ID), filmID, cfg, repo)

				// if err != nil {
				// 	fmt.Printf("❌ Failed to download film %d (ID: %d): %v\n", film.ID, film.ID, err)

				// 	continue
				// }
			}

		}

		fmt.Printf("✅ Successfully processed film: %d\n", film.ID)
	}

	return nil
}
