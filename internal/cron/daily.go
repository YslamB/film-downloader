package cron

import (
	"context"
	"film-downloader/internal/config"
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
	fmt.Println("asdfoiufjkio")
	var movies []models.Movie
	var err error

	if episodeID == "" && seasonID == "" && filmID != "" {
		movieSources, err := DownloadMovieSourceWithID(ctx, filmID, cfg, repo)

		if err != nil {
			fmt.Println("soidhfi")
			return err
		}
		movies = append(movies, movieSources...)
	}

	if seasonID != "" {
		fmt.Println("🔍 Checking season with ID:", seasonID)
		movies, err = requests.GetEpisodesSourceWithSeasonID(seasonID, episodeID, cfg)
		time.Sleep(1 * time.Second)
		fmt.Println("s89dhuinuj")

		if err != nil {
			fmt.Println("sd89fhuin")
			return err
		}
	}
	fmt.Println("so9d8fuhin")

	fmt.Println("✅ Received Source files...", movies)

	for i := range movies {
		// err := downloader.DownloadHLS(movies[i], cfg)

		// if err != nil {
		// 	return err
		// }

		// err = utils.UploadFolderToMinio(
		// 	"temp/"+movies[i].Name, movies[i].Name, cfg.MINIO_BUCKET,
		// 	cfg.MINIO_ENDPOINT, cfg.MINIO_ACCESS_KEY, cfg.MINIO_SECRET_KEY,
		// 	cfg.MINIO_SECURE, cfg.MINIO_WORKERS,
		// )

		// if err != nil {
		// 	return err
		// }

		fmt.Println("soidfhu8i9")
		fileID, err := repo.GetFileID(ctx, movies[i].Name)

		if err != nil {
			fmt.Println("09s8duhinj")
			return err
		}

		err = repo.CreateMovieFile(ctx, fileID, movies[i].ID)

		if err != nil {
			fmt.Println("asdfoiu9u83ybhwsfjkio")
			return err
		}

		fmt.Println("🔍 File ID:", fileID)
		err = os.RemoveAll("temp/" + movies[i].Name)

		if err != nil {
			fmt.Println("893uhienjf")
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

	for i := range searchResult.Films {

		filmID := fmt.Sprintf("%d", searchResult.Films[i].ID)

		if searchResult.Films[i].TypeID == 1 {
			err := DownloadWithID(ctx, "", "", filmID, cfg, repo)

			if err != nil {
				fmt.Println("❌ Failed to download film", filmID, err)
				continue
			}
		} else {
			continue
			// seasons, err := requests.GetSeasonsData(ctx, filmID, cfg)

			// if err != nil {
			// 	fmt.Println("❌ Failed to get seasons data for film", filmID, err)
			// 	continue
			// }

			// for i := range seasons {
			// 	err := DownloadWithID(ctx, "", fmt.Sprintf("%d", seasons[i].ID), filmID, cfg, repo)

			// 	if err != nil {
			// 		fmt.Println("❌ Failed to download season", seasons[i].ID, err)
			// 		continue
			// 	}
			// }

		}

	}

	return nil
}
