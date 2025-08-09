package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetFilmSourceURL(filmID string) (string, error) {
	apiURL := fmt.Sprintf("https://film.beletapis.com/api/v2/files/%s?type=1", filmID)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad response: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var result struct {
		Sources []struct {
			Filename string `json:"filename"`
			Quality  string `json:"quality"`
		} `json:"sources"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}

	for _, source := range result.Sources {
		if source.Quality == "1080p" {
			return source.Filename, nil
		}
	}

	return "", fmt.Errorf("1080p source not found")
}
