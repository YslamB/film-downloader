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

func DownloadWithID(ctx context.Context, episodeID string, season *models.Season, filmID string, cfg *config.Config, repo *repositories.MovieRepository) error {
	var movies []models.Movie
	var err error

	if season == nil && filmID != "" {
		movieSources, err := DownloadMovieSourceWithID(ctx, filmID, cfg, repo)

		if err != nil {
			return err
		}
		movies = append(movies, movieSources...)

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

			fileID, err := repo.GetFileID(ctx, movies[i].Name)

			if err != nil {
				return err
			}

			err = repo.CreateMovieFile(ctx, fileID, movies[i].ID)

			if err != nil {
				return err
			}

			err = os.RemoveAll("temp/" + movies[i].Name)

			if err != nil {
				return err
			}
		}

	}

	if season != nil {
		movies, err = requests.GetEpisodesSourceWithSeasonID(ctx, season, cfg, repo)
		time.Sleep(1 * time.Second)

		if err != nil {
			return err
		}
	}

	return nil
}

func GetLastMovies(ctx context.Context, cfg *config.Config, repo *repositories.MovieRepository, wg *sync.WaitGroup) error {
	searchResult, err := requests.GetSearchResults(ctx, 1, cfg)

	if err != nil {
		return fmt.Errorf("failed to get search results from API: %w", err)
	}

	for i := range searchResult.Films {
		wg.Add(1)

		filmID := fmt.Sprintf("%d", searchResult.Films[i].ID)

		if searchResult.Films[i].TypeID == 1 {
			err := DownloadWithID(ctx, "", nil, filmID, cfg, repo)
			wg.Done()

			if err != nil {
				fmt.Println("Error downloading movie:", err)
				continue
			}

		} else {
			bbmovieID, err := CreateMovie(ctx, filmID, cfg, repo)

			if err != nil {
				fmt.Println("Error creating movie:", err)
			}

			seasons, err := CreateMovieSeasons(ctx, fmt.Sprintf("%d", bbmovieID), filmID, cfg, repo)

			if err != nil {
				fmt.Println("Error creating movie seasons:", err)
				continue
			}

			for i := range seasons {
				fmt.Println("IIIIII  :::  ", i)
				err := DownloadWithID(ctx, "", &seasons[i], filmID, cfg, repo)

				if err != nil {
					fmt.Println("Error downloading movie episodes:", err)
					continue
				}
			}

			wg.Done()
			continue
			// seasons, err := requests.GetSeasonsData(ctx, filmID, cfg)

			// if err != nil {
			// 	continue
			// }

			// for i := range seasons {
			// 	err := DownloadWithID(ctx, "", fmt.Sprintf("%d", seasons[i].ID), filmID, cfg, repo)

			// 	if err != nil {
			// 		continue
			// 	}
			// }

		}

	}

	return nil
}

func RefreshToken(ctx context.Context, cfg *config.Config, repo *repositories.MovieRepository) error {
	err := repo.RefreshToken(ctx)

	if err != nil {
		return fmt.Errorf("failed to refresh token: %w", err)
	}

	return nil
}
