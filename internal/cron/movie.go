package cron

import (
	"context"
	"film-downloader/internal/config"
	"film-downloader/internal/models"
	"film-downloader/internal/repositories"
	"film-downloader/internal/requests"
	"fmt"
	"strconv"
	"time"
)

func DownloadMovieSourceWithID(ctx context.Context, filmID string, cfg *config.Config, repo *repositories.MovieRepository) (models.Movie, error) {
	var movieSource models.Movie
	movieRes, err := requests.GetMovieData(ctx, filmID, cfg)

	if err != nil {
		return movieSource, err
	}

	exists, err := repo.CheckMovieExists(ctx, strconv.Itoa(movieRes.Film.ID))
	time.Sleep(1 * time.Second)

	if err != nil {
		return movieSource, err
	}

	if exists {
		return movieSource, fmt.Errorf("movieRes already exists")
	}

	movieRes.Film.CategoryID, err = repo.GetCategoryID(ctx, movieRes.Film.CategoryID)
	fmt.Println("🔍 Category ID:", movieRes.Film.CategoryID)
	time.Sleep(1 * time.Second)

	if err != nil {
		return movieSource, err
	}

	genreIDs, err := repo.GetGenreIDs(ctx, movieRes.Film.Genres)
	fmt.Println("🔍 Genre IDs:", genreIDs)
	time.Sleep(1 * time.Second)

	if err != nil {
		return movieSource, err
	}

	countryIDs, err := repo.GetCountryIDs(ctx, movieRes.Film.Countries)
	fmt.Println("🔍 Country IDs:", countryIDs)
	time.Sleep(1 * time.Second)

	if err != nil {
		return movieSource, err
	}

	actorIDs, err := repo.GetActorIDs(ctx, movieRes.Film.Actors)
	if err != nil {
		return movieSource, err
	}
	fmt.Println("🔍 Actor IDs:", actorIDs)
	time.Sleep(1 * time.Second)

	directorIDs, err := repo.GetActorIDs(ctx, movieRes.Film.Directors)
	if err != nil {
		return movieSource, err
	}
	fmt.Println("🔍 Director IDs:", directorIDs)
	time.Sleep(1 * time.Second)

	studioIDs, err := repo.GetStudioIDs(ctx, movieRes.Film.Studios)
	if err != nil {
		return movieSource, err
	}
	fmt.Println("🔍 Studio IDs:", studioIDs)
	time.Sleep(1 * time.Second)

	verticalImageID, err := repo.SendMovieImage(ctx, movieRes.Film.Images.Vertical.Default)
	if err != nil {
		return movieSource, err
	}
	fmt.Println("🔍 Vertical Image ID:", verticalImageID)
	time.Sleep(1 * time.Second)

	verticalWithoutNameImageID, err := repo.SendMovieImage(ctx, movieRes.Film.Images.VerticalWithoutName.Default)
	if err != nil {
		return movieSource, err
	}
	fmt.Println("🔍 Vertical Without Name Image ID:", verticalWithoutNameImageID)
	time.Sleep(1 * time.Second)

	horizontalWithNameImageID, err := repo.SendMovieImage(ctx, movieRes.Film.Images.HorizontalWithName.Default)
	if err != nil {
		return movieSource, err
	}
	fmt.Println("🔍 Horizontal With Name Image ID:", horizontalWithNameImageID)
	time.Sleep(1 * time.Second)

	horizontalWithoutNameImageID, err := repo.SendMovieImage(ctx, movieRes.Film.Images.HorizontalWithoutName.Default)
	if err != nil {
		return movieSource, err
	}
	fmt.Println("🔍 Horizontal Without Name Image ID:", horizontalWithoutNameImageID)
	time.Sleep(1 * time.Second)

	nameImageID, err := repo.SendMovieImage(ctx, movieRes.Film.Images.Name)
	if err != nil {
		return movieSource, err
	}
	fmt.Println("🔍 Name Image ID:", nameImageID)
	time.Sleep(1 * time.Second)

	languageID, err := repo.GetLanguageID(ctx, movieRes.Film.Language)
	if err != nil {
		return movieSource, err
	}
	fmt.Println("🔍 Language ID:", languageID)
	time.Sleep(1 * time.Second)

	id, err := repo.CreateMovie(ctx, movieRes.Film, genreIDs, countryIDs, actorIDs, directorIDs, studioIDs, languageID, verticalImageID, verticalWithoutNameImageID, horizontalWithNameImageID, horizontalWithoutNameImageID, nameImageID)

	if err != nil {
		return movieSource, err
	}
	fmt.Println("🔍 Movie ID:", id)
	time.Sleep(1 * time.Second)

	movieSource, err = requests.GetFilmSourceURL(ctx, filmID, cfg)

	if err != nil {
		return movieSource, err
	}

	return movieSource, nil
}
