package requests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"film-downloader/internal/config"
	"film-downloader/internal/models"
)

func GetMovieData(ctx context.Context, movieID string, cfg *config.Config) (models.MovieResponse, error) {
	var movieResponse models.MovieResponse

	// Construct the API URL
	url := fmt.Sprintf("https://film.beletapis.com/api/v2/movie/%s", movieID)

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return movieResponse, fmt.Errorf("failed to create request: %w", err)
	}

	// Add Chrome-like headers
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

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return movieResponse, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return movieResponse, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	// Decode JSON response
	if err := json.NewDecoder(resp.Body).Decode(&movieResponse); err != nil {
		return movieResponse, fmt.Errorf("failed to decode response: %w", err)
	}

	return movieResponse, nil
}
