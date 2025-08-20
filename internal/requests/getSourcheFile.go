package requests

import (
	"context"
	"encoding/json"
	"film-downloader/internal/config"
	"film-downloader/internal/models"
	"film-downloader/internal/utils"
	"fmt"
	"io"
	"net/http"
	"time"
)

func GetFilmSourceURL(ctx context.Context, filmID string, cfg *config.Config) (models.Movie, error) {
	var movie models.Movie
	apiURL := fmt.Sprintf("https://film.beletapis.com/api/v2/files/%s?type=1", filmID)
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)

	if err != nil {
		return movie, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", cfg.AccessToken)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Sec-Ch-Ua", `"Google Chrome";v="119", "Chromium";v="119", "Not?A_Brand";v="24"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"macOS"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")

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
			movie.Name, err = utils.GenerateUUID()

			if err != nil {
				return movie, fmt.Errorf("failed to generate UUID: %w", err)
			}
			return movie, nil
		}
	}

	return movie, fmt.Errorf("1080p source not found")
}
