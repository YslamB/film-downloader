package cron

import (
	"context"
	"errors"
	"film-downloader/internal/config"
	"film-downloader/internal/models"
	"film-downloader/internal/repositories"
	"film-downloader/internal/requests"
	"strconv"
	"time"
)

func DownloadMovieSourceWithID(ctx context.Context, filmID string, cfg *config.Config, repo *repositories.MovieRepository) ([]models.Movie, error) {
	var movieSources []models.Movie
	movieID, err := CreateMovie(ctx, filmID, cfg, repo)

	if err != nil {
		return movieSources, err
	}

	movieSources, err = requests.GetFilmSourceURL(ctx, filmID, cfg, movieID)

	if err != nil {
		return movieSources, err
	}

	return movieSources, nil
}

func CreateMovie(ctx context.Context, filmID string, cfg *config.Config, repo *repositories.MovieRepository) (int, error) {
	movieRes, err := requests.GetMovieData(ctx, filmID, cfg)

	if err != nil {
		return 0, err
	}

	movieID, err := repo.CheckMovieExists(ctx, strconv.Itoa(movieRes.Film.ID))

	if err != nil {
		return 0, err
	}

	if movieID != 0 {
		return movieID, errors.New("movie already exists")
	}

	movieRes.Film.CategoryID, err = repo.GetCategoryID(ctx, movieRes.Film.CategoryID)
	time.Sleep(1 * time.Second)

	if err != nil {
		return 0, err
	}

	genreIDs, err := repo.GetGenreIDs(ctx, movieRes.Film.Genres)
	time.Sleep(1 * time.Second)

	if err != nil {
		return 0, err
	}

	countryIDs, err := repo.GetCountryIDs(ctx, movieRes.Film.Countries)
	time.Sleep(1 * time.Second)

	if err != nil {
		return 0, err
	}

	actorIDs, err := repo.GetActorIDs(ctx, movieRes.Film.Actors)

	if err != nil {
		return 0, err
	}

	time.Sleep(1 * time.Second)
	directorIDs, err := repo.GetActorIDs(ctx, movieRes.Film.Directors)

	if err != nil {
		return 0, err
	}

	time.Sleep(1 * time.Second)
	studioIDs, err := repo.GetStudioIDs(ctx, movieRes.Film.Studios)

	if err != nil {
		return 0, err
	}

	time.Sleep(1 * time.Second)
	verticalImageID, err := repo.SendMovieImage(ctx, movieRes.Film.Images.Vertical.Default)

	if err != nil {
		return 0, err
	}

	time.Sleep(1 * time.Second)
	verticalWithoutNameImageID, err := repo.SendMovieImage(ctx, movieRes.Film.Images.VerticalWithoutName.Default)

	if err != nil {
		return 0, err
	}

	time.Sleep(1 * time.Second)
	horizontalWithNameImageID, err := repo.SendMovieImage(ctx, movieRes.Film.Images.HorizontalWithName.Default)

	if err != nil {
		return 0, err
	}

	time.Sleep(1 * time.Second)
	horizontalWithoutNameImageID, err := repo.SendMovieImage(ctx, movieRes.Film.Images.HorizontalWithoutName.Default)

	if err != nil {
		return 0, err
	}

	time.Sleep(1 * time.Second)
	nameImageID, err := repo.SendMovieImage(ctx, movieRes.Film.Images.Name)

	if err != nil {
		return 0, err
	}

	time.Sleep(1 * time.Second)
	languageID, err := repo.GetLanguageID(ctx, movieRes.Film.Language)

	if err != nil {
		return 0, err
	}

	time.Sleep(1 * time.Second)
	movieID, err = repo.CreateMovie(
		ctx, movieRes.Film, genreIDs, countryIDs, actorIDs, directorIDs,
		studioIDs, languageID, verticalImageID, verticalWithoutNameImageID,
		horizontalWithNameImageID, horizontalWithoutNameImageID, nameImageID,
	)

	if err != nil {
		return 0, err
	}

	time.Sleep(1 * time.Second)
	return movieID, nil
}

func CreateMovieSeasons(ctx context.Context, movieID, filmID string, cfg *config.Config, repo *repositories.MovieRepository) ([]models.Season, error) {
	movieRes, err := requests.GetMovieData(ctx, filmID, cfg)

	if err != nil {
		return nil, err
	}

	for i := range movieRes.Film.Seasons {
		seasonID, err := repo.CreateSeason(ctx, movieRes.Film.Seasons[i], movieID)

		if err != nil {
			return movieRes.Film.Seasons, err
		}

		movieRes.Film.Seasons[i].BBID = seasonID
	}

	return movieRes.Film.Seasons, nil
}
