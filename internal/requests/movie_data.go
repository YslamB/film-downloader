package requests

import (
	"context"
	"fmt"
	"time"

	"film-downloader/internal/config"
	"film-downloader/internal/models"
	"film-downloader/internal/utils"
)

const (
	movieDataURL = "https://film.beletapis.com/api/v2/movie/%s"
	searchURL    = "https://film-search.belet.me/api/v1/search"
)

func GetMovieData(ctx context.Context, movieID string, cfg *config.Config) (models.MovieResponse, error) {
	var movieResponse models.MovieResponse

	apiConfig := utils.APIRequestConfig{
		Method:      "GET",
		URL:         fmt.Sprintf(movieDataURL, movieID),
		AccessToken: cfg.GetAccessToken(),
		Timeout:     30 * time.Second,
	}

	err := utils.MakeJSONRequest(ctx, apiConfig, &movieResponse)

	if err != nil {
		return movieResponse, utils.WrapErrorf(err, "failed to get movie data for ID %s", movieID)
	}

	return movieResponse, nil
}

func GetSearchResults(ctx context.Context, page int, cfg *config.Config) (models.SearchResult, error) {
	var searchResult models.SearchResult

	requestBody := models.SearchRequest{
		Page:  page,
		Order: "desc",
		Sort:  "add",
	}

	apiConfig := utils.APIRequestConfig{
		Method:      "POST",
		URL:         searchURL,
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
