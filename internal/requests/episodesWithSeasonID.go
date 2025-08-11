package requests

import (
	"encoding/json"
	"film-downloader/internal/config"
	"film-downloader/internal/models"
	"fmt"
	"io"
	"net/http"
)

func GetEpisodesWithSeasonID(seasonID, episodeID string, cfg config.Config) ([]models.Movie, error) {
	var movies []models.Movie
	url := fmt.Sprintf("https://film.beletapis.com/api/v2/episodes?seasonId=%s", seasonID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return movies, fmt.Errorf("❌ failed to create request: %w", err)
	}
	req.Header.Set("Authorization", cfg.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return movies, fmt.Errorf("❌ request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return movies, fmt.Errorf("❌ bad response: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return movies, fmt.Errorf("❌ failed to read response: %w", err)
	}

	var result struct {
		Episodes []struct {
			ID      int    `json:"id"`
			Name    string `json:"name"`
			Sources []struct {
				DownloadURL string `json:"download_url"`
				Quality     string `json:"quality"`
			} `json:"sources"`
		} `json:"episodes"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return movies, fmt.Errorf("❌ failed to parse JSON: %w", err)
	}

	for _, ep := range result.Episodes {
		idStr := fmt.Sprintf("%d", ep.ID)
		for _, source := range ep.Sources {
			if source.Quality == "1080p" {
				if episodeID != "" && idStr == episodeID {
					movies = append(movies, models.Movie{Source: source.DownloadURL, Name: ep.Name})
					return movies, nil
				}
				if episodeID == "" {
					movies = append(movies, models.Movie{Source: source.DownloadURL, Name: ep.Name})
				}
			}
		}
	}

	// If episodeID was specified but not found
	if episodeID != "" && len(movies) == 0 {
		return movies, fmt.Errorf("❌ episode %s not found or has no 1080p source", episodeID)
	}

	return movies, nil
}
