package requests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"film-downloader/internal/config"
	"film-downloader/internal/models"
)

func GetMovieData(ctx context.Context, movieID string, cfg *config.Config) (models.MovieResponse, error) {
	var movieResponse models.MovieResponse

	url := fmt.Sprintf("https://film.beletapis.com/api/v2/movie/%s", movieID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return movieResponse, fmt.Errorf("failed to create request: %w", err)
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
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return movieResponse, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return movieResponse, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&movieResponse); err != nil {
		return movieResponse, fmt.Errorf("failed to decode response: %w", err)
	}

	return movieResponse, nil
}

func GetSeasonsData(ctx context.Context, movieID string, cfg *config.Config) ([]models.Season, error) {
	var movieResponse models.MovieResponse
	url := fmt.Sprintf("https://film.beletapis.com/api/v2/movie/%s", movieID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
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
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&movieResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return movieResponse.Film.Seasons, nil
}

// SearchRequest represents the request body for the search API
type SearchRequest struct {
	Page  int    `json:"page"`
	Order string `json:"order"`
}

// GetSearchResults sends a POST request to the search API and returns the search results
func GetSearchResults(ctx context.Context, page int, cfg *config.Config) (models.SearchResult, error) {
	var searchResult models.SearchResult

	url := "https://film-search.belet.me/api/v1/search"

	// Create request body
	requestBody := SearchRequest{
		Page:  page,
		Order: "desc",
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return searchResult, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return searchResult, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", cfg.AccessToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)

	if err != nil {
		return searchResult, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("API Error - Status: %d, Body: %s\n", resp.StatusCode, string(body))
		return searchResult, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&searchResult); err != nil {
		return searchResult, fmt.Errorf("failed to decode response: %w", err)
	}

	return searchResult, nil
}
