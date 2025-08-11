package requests

import (
	"encoding/json"
	"film-downloader/internal/config"
	"film-downloader/internal/models"
	"fmt"
	"io"
	"net/http"
	"time"
)

func GetFilmSourceURL(filmID string, cfg config.Config) (models.Movie, error) {
	var movie models.Movie
	apiURL := fmt.Sprintf("https://film.beletapis.com/api/v2/files/%s?type=1", filmID)

	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return movie, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", cfg.AccessToken)

	client := &http.Client{
		Timeout: time.Second * 5,
	}
	resp, err := client.Do(req)
	if err != nil {
		return movie, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return movie, fmt.Errorf("bad response: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return movie, fmt.Errorf("failed to read response: %w", err)
	}

	var result struct {
		Sources []struct {
			Filename string `json:"filename"`
			Quality  string `json:"quality"`
		} `json:"sources"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return movie, fmt.Errorf("failed to parse JSON: %w", err)
	}

	for _, source := range result.Sources {
		if source.Quality == "1080p" {
			movie.Source = source.Filename
			movie.Name = "1080p"
			return movie, nil
		}
	}

	return movie, fmt.Errorf("1080p source not found")
}
