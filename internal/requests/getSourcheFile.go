package requests

import (
	"context"
	"film-downloader/internal/config"
	"film-downloader/internal/models"
	"film-downloader/internal/utils"
	"fmt"
	"time"
)

func GetFilmSourceURL(ctx context.Context, filmID string, cfg *config.Config, movieID int) ([]models.Movie, error) {
	var movie models.Movie
	var movies []models.Movie
	apiURL := fmt.Sprintf("https://film.beletapis.com/api/v2/files/%s?type=1", filmID)

	var result struct {
		Sources []struct {
			Filename string `json:"filename"`
			Quality  string `json:"quality"`
			Type     string `json:"type"`
		} `json:"sources"`
	}

	apiConfig := utils.APIRequestConfig{
		Method:      "GET",
		URL:         apiURL,
		AccessToken: cfg.GetAccessToken(),
		Timeout:     5 * time.Second,
	}

	err := utils.MakeJSONRequest(ctx, apiConfig, &result)

	if err != nil {
		return []models.Movie{}, utils.WrapErrorf(err, "failed to get film source URL for film ID %s", filmID)
	}

	for _, source := range result.Sources {
		main := source.Quality == "1080p"

		movie.Sources = append(movie.Sources, models.Source{
			MasterFile: source.Filename,
			Quality:    source.Quality,
			Type:       source.Type,
			Main:       main,
		})
	}

	movie.Name, err = utils.GenerateUUID()

	if err != nil {
		return []models.Movie{}, utils.WrapError(err, "failed to generate UUID")
	}

	movie.ID = movieID
	movies = append(movies, movie)

	return movies, nil
}
