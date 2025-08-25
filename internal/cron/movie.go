package cron

import (
	"context"
	"film-downloader/internal/config"
	"film-downloader/internal/models"
	"film-downloader/internal/repositories"
	"film-downloader/internal/requests"
	"film-downloader/internal/utils"
	"fmt"
	"strconv"
)

func DownloadMovieSourceWithID(ctx context.Context, filmID string, cfg *config.Config, repo *repositories.MovieRepository) ([]models.Movie, error) {
	var movieSources []models.Movie
	delayManager := utils.NewDelayManager()
	delayManager.SetAPIDelay()
	movieRes, err := requests.GetMovieData(ctx, filmID, cfg)

	if err != nil {
		return movieSources, err
	}

	exists, err := repo.CheckMovieExists(ctx, strconv.Itoa(movieRes.Film.ID))
	delayManager.Execute("api")

	if err != nil {
		return movieSources, err
	}

	if exists {
		return movieSources, fmt.Errorf("movieRes already exists")
	}

	movieRes.Film.CategoryID, err = repo.GetCategoryID(ctx, movieRes.Film.CategoryID)
	fmt.Println("🔍 Category ID:", movieRes.Film.CategoryID)
	delayManager.Execute("api")

	if err != nil {
		return movieSources, err
	}

	genreIDs, err := repo.GetGenreIDs(ctx, movieRes.Film.Genres)
	fmt.Println("🔍 Genre IDs:", genreIDs)
	delayManager.Execute("api")

	if err != nil {
		return movieSources, err
	}

	countryIDs, err := repo.GetCountryIDs(ctx, movieRes.Film.Countries)
	fmt.Println("🔍 Country IDs:", countryIDs)
	delayManager.Execute("api")

	if err != nil {
		return movieSources, err
	}

	actorIDs, err := repo.GetActorIDs(ctx, movieRes.Film.Actors)

	if err != nil {
		return movieSources, err
	}

	fmt.Println("🔍 Actor IDs:", actorIDs)
	delayManager.Execute("api")
	directorIDs, err := repo.GetActorIDs(ctx, movieRes.Film.Directors)

	if err != nil {
		return movieSources, err
	}

	fmt.Println("🔍 Director IDs:", directorIDs)
	delayManager.Execute("api")
	studioIDs, err := repo.GetStudioIDs(ctx, movieRes.Film.Studios)

	if err != nil {
		return movieSources, err
	}

	fmt.Println("🔍 Studio IDs:", studioIDs)
	delayManager.Execute("api")
	verticalImageID, err := repo.SendMovieImage(ctx, movieRes.Film.Images.Vertical.Default)

	if err != nil {
		return movieSources, err
	}

	fmt.Println("🔍 Vertical Image ID:", verticalImageID)
	delayManager.Execute("api")
	verticalWithoutNameImageID, err := repo.SendMovieImage(ctx, movieRes.Film.Images.VerticalWithoutName.Default)

	if err != nil {
		return movieSources, err
	}

	fmt.Println("🔍 Vertical Without Name Image ID:", verticalWithoutNameImageID)
	delayManager.Execute("api")
	horizontalWithNameImageID, err := repo.SendMovieImage(ctx, movieRes.Film.Images.HorizontalWithName.Default)

	if err != nil {
		return movieSources, err
	}

	fmt.Println("🔍 Horizontal With Name Image ID:", horizontalWithNameImageID)
	delayManager.Execute("api")
	horizontalWithoutNameImageID, err := repo.SendMovieImage(ctx, movieRes.Film.Images.HorizontalWithoutName.Default)

	if err != nil {
		return movieSources, err
	}

	fmt.Println("🔍 Horizontal Without Name Image ID:", horizontalWithoutNameImageID)
	delayManager.Execute("api")
	nameImageID, err := repo.SendMovieImage(ctx, movieRes.Film.Images.Name)

	if err != nil {
		return movieSources, err
	}

	fmt.Println("🔍 Name Image ID:", nameImageID)
	delayManager.Execute("api")
	languageID, err := repo.GetLanguageID(ctx, movieRes.Film.Language)

	if err != nil {
		return movieSources, err
	}

	fmt.Println("🔍 Language ID:", languageID)
	delayManager.Execute("api")
	movieID, err := repo.CreateMovie(
		ctx, movieRes.Film, genreIDs, countryIDs, actorIDs, directorIDs,
		studioIDs, languageID, verticalImageID, verticalWithoutNameImageID,
		horizontalWithNameImageID, horizontalWithoutNameImageID, nameImageID,
	)

	if err != nil {
		return movieSources, err
	}

	fmt.Println("🔍 Movie ID:", movieID)
	delayManager.Execute("api")
	movieSources, err = requests.GetFilmSourceURL(ctx, filmID, cfg, movieID)

	if err != nil {
		return movieSources, err
	}

	return movieSources, nil
}
