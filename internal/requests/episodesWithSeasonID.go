package requests

import (
	"context"
	"encoding/json"
	"film-downloader/internal/config"
	"film-downloader/internal/models"
	"film-downloader/internal/repositories"
	"film-downloader/internal/utils"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	episodesURL = "https://film.beletapis.com/api/v2/episodes?seasonId=%s"
)

func GetEpisodesSourceWithSeasonID(ctx context.Context, season *models.Season, cfg *config.Config, repo *repositories.MovieRepository) ([]models.Movie, error) {
	var movies []models.Movie
	url := fmt.Sprintf(episodesURL, fmt.Sprintf("%d", season.ID))
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return movies, fmt.Errorf("❌ failed to create request: %w", err)
	}

	req.Header.Set("Authorization", cfg.GetAccessToken())
	utils.SetCommonHeaders(req, cfg.GetAccessToken())

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

	var result models.EpisodeResponse

	if err := json.Unmarshal(body, &result); err != nil {
		return movies, fmt.Errorf("❌ failed to parse JSON: %w", err)
	}

	for _, ep := range result.Episodes {
		// create episode
		ep.FilePath, err = utils.GenerateUUID()

		if err != nil {
			return movies, fmt.Errorf("❌ failed to generate UUID: %w", err)
		}

		ep.FilePath += fmt.Sprintf("/seasons/%d", season.BBID)
		ep.FileID, err = repo.GetFileID(ctx, ep.FilePath)

		if err != nil {
			return movies, fmt.Errorf("❌ failed to get file ID: %w", err)
		}

		episodeBBID, err := repo.CreateEpisode(ctx, ep, season.BBID)
		time.Sleep(1 * time.Second)

		if err != nil {
			return movies, fmt.Errorf("❌ failed to create episode: %w", err)
		}

		movie := models.Movie{Name: ep.FilePath, ID: episodeBBID, Type: models.EpisodeType}

		for i := range ep.Sources {
			main := false

			if ep.Sources[i].Quality == "1080p" {
				main = true
			}

			movie.Sources = append(movie.Sources, models.Source{
				MasterFile: ep.Sources[i].DownloadURL,
				Quality:    ep.Sources[i].Quality,
				Main:       main,
			})

			movies = append(movies, movie)
		}
	}

	if season != nil && len(movies) == 0 {
		return movies, fmt.Errorf("❌ season %s not found or has no 1080p source", season.Name)
	}

	return movies, nil
}
