package requests

import (
	"context"
	"fmt"
	"time"

	"film-downloader/internal/config"
	"film-downloader/internal/models"
	"film-downloader/internal/utils"
)

func GetMovieData(ctx context.Context, movieID string, cfg *config.Config) (models.MovieResponse, error) {
	var movieResponse models.MovieResponse

	url := fmt.Sprintf("https://film.beletapis.com/api/v2/movie/%s", movieID)

	apiConfig := utils.APIRequestConfig{
		Method:      "GET",
		URL:         url,
		AccessToken: cfg.GetAccessToken(),
		Timeout:     30 * time.Second,
	}

	err := utils.MakeJSONRequest(ctx, apiConfig, &movieResponse)
	if err != nil {
		return movieResponse, utils.WrapErrorf(err, "failed to get movie data for ID %s", movieID)
	}

	return movieResponse, nil
}

func GetSeasonsData(ctx context.Context, movieID string, cfg *config.Config) ([]models.Season, error) {
	var movieResponse models.MovieResponse
	url := fmt.Sprintf("https://film.beletapis.com/api/v2/movie/%s", movieID)

	apiConfig := utils.APIRequestConfig{
		Method:      "GET",
		URL:         url,
		AccessToken: cfg.GetAccessToken(),
		Timeout:     10 * time.Second,
	}

	err := utils.MakeJSONRequest(ctx, apiConfig, &movieResponse)
	if err != nil {
		return nil, utils.WrapErrorf(err, "failed to get seasons data for movie ID %s", movieID)
	}

	return movieResponse.Film.Seasons, nil
}

type SearchRequest struct {
	Page  int    `json:"page"`
	Order string `json:"order"`
}

func GetSearchResults(ctx context.Context, page int, cfg *config.Config) (models.SearchResult, error) {
	var searchResult models.SearchResult

	url := "https://film-search.belet.me/api/v1/search"

	requestBody := SearchRequest{
		Page:  page,
		Order: "desc",
	}

	apiConfig := utils.APIRequestConfig{
		Method:      "POST",
		URL:         url,
		Body:        requestBody,
		AccessToken: cfg.GetAccessToken(),
		Timeout:     30 * time.Second,
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
	}

	err := utils.MakeJSONRequest(ctx, apiConfig, &searchResult)
	if err != nil {
		return searchResult, utils.WrapErrorf(err, "failed to get search results for page %d", page)
	}

	return searchResult, nil
}
