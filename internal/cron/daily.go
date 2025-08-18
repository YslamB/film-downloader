package cron

import (
	"context"
	"film-downloader/internal/config"
	"film-downloader/internal/downloader"
	"film-downloader/internal/models"
	"film-downloader/internal/repositories"
	"film-downloader/internal/requests"
	"fmt"
	"strconv"
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
		movieRes, err := requests.GetMovieData(ctx, filmID, cfg)

		if err != nil {
			return err
		}

		exists, err := repo.CheckMovieExists(ctx, strconv.Itoa(movieRes.Film.ID))

		if err != nil {
			return err
		}

		if exists {
			return fmt.Errorf("movieRes already exists")
		}

		movieRes.Film.CategoryID, err = repo.GetCategoryID(ctx, movieRes.Film.CategoryID)

		if err != nil {
			return err
		}

		genreIDs, err := repo.GetGenreIDs(ctx, movieRes.Film.Genres)

		if err != nil {
			return err
		}

		countryIDs, err := repo.GetCountryIDs(ctx, movieRes.Film.Countries)

		if err != nil {
			return err
		}

		actorIDs, err := repo.GetActorIDs(ctx, movieRes.Film.Actors)

		if err != nil {
			return err
		}

		err = repo.CreateMovie(ctx, movieRes, genreIDs, countryIDs, actorIDs)

		if err != nil {
			return err
		}

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
