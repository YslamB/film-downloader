package requests

import (
	"encoding/json"
	"film-downloader/internal/config"
	"film-downloader/internal/models"
	"film-downloader/internal/utils"
	"fmt"
	"io"
	"net/http"
)

func GetEpisodesSourceWithSeasonID(seasonID, episodeID string, cfg *config.Config) ([]models.Movie, error) {
	var movies []models.Movie
	url := fmt.Sprintf("https://film.beletapis.com/api/v2/episodes?seasonId=%s", seasonID)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return movies, fmt.Errorf("❌ failed to create request: %w", err)
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
			main := false

			if source.Quality == "1080p" {
				main = true
			}

			uuid, err := utils.GenerateUUID()

			if err != nil {
				return movies, fmt.Errorf("failed to generate UUID: %w", err)
			}

			movie := models.Movie{Name: uuid}

			if episodeID != "" && idStr == episodeID {
				movie.Sources = append(movie.Sources, models.Source{
					MasterFile: source.DownloadURL,
					Quality:    source.Quality,
					Main:       main,
				})
			}

			if episodeID == "" {
				movie.Sources = append(movie.Sources, models.Source{
					MasterFile: source.DownloadURL,
					Quality:    source.Quality,
					Main:       main,
				})
			}

			movies = append(movies, movie)
		}
	}

	if episodeID != "" && len(movies) == 0 {
		return movies, fmt.Errorf("❌ episode %s not found or has no 1080p source", episodeID)
	}

	return movies, nil
}
