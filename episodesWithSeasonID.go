package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetEpisodesWithSeasonID(seasonID, episodeID string) ([]string, error) {
	url := fmt.Sprintf("https://film.beletapis.com/api/v2/episodes?seasonId=%s", seasonID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("❌ failed to create request: %w", err)
	}
	req.Header.Set("Authorization", accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("❌ request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("❌ bad response: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("❌ failed to read response: %w", err)
	}

	var result struct {
		Episodes []struct {
			ID      int `json:"id"`
			Sources []struct {
				DownloadURL string `json:"download_url"`
				Quality     string `json:"quality"`
			} `json:"sources"`
		} `json:"episodes"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("❌ failed to parse JSON: %w", err)
	}

	var urls []string
	for _, ep := range result.Episodes {
		idStr := fmt.Sprintf("%d", ep.ID)
		for _, source := range ep.Sources {
			if source.Quality == "1080p" {
				if episodeID != "" && idStr == episodeID {
					return []string{source.DownloadURL}, nil
				}
				if episodeID == "" {
					urls = append(urls, source.DownloadURL)
				}
			}
		}
	}

	// If episodeID was specified but not found
	if episodeID != "" && len(urls) == 0 {
		return nil, fmt.Errorf("❌ episode %s not found or has no 1080p source", episodeID)
	}

	return urls, nil
}
