package repositories

import (
	"bytes"
	"context"
	"encoding/json"
	"film-downloader/internal/models"
	"fmt"
	"net/http"
	"time"
)

type MovieRepository struct {
	accessToken string
}

const (
	getMovieURL      = "http://95.85.126.217:5050/api/v1/admin/movies/%s"
	getCategoryIDURL = "http://95.85.126.217:5050/api/v1/admin/catalogs/categories/%d"
	getGenreIDURL    = "http://95.85.126.217:5050/api/v1/admin/catalogs/genres"
	getCountryIDURL  = "http://95.85.126.217:5050/api/v1/admin/catalogs/countries"
	getActorIDURL    = "http://95.85.126.217:5050/api/v1/admin/catalogs/persons"
)

func NewMovieRepository(accessToken string) *MovieRepository {
	return &MovieRepository{
		accessToken: accessToken,
	}
}

func (r *MovieRepository) CheckMovieExists(ctx context.Context, movieID string) (bool, error) {
	// req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(getMovieURL, movieID), nil)

	// if err != nil {
	// 	return false, fmt.Errorf("failed to create request: %w", err)
	// }

	// req.Header.Set("Authorization", r.accessToken)

	// client := &http.Client{
	// 	Timeout: time.Second * 10,
	// }
	// resp, err := client.Do(req)

	// if err != nil {
	// 	return false, fmt.Errorf("request failed: %w", err)
	// }

	// defer resp.Body.Close()

	// if resp.StatusCode != http.StatusOK {
	// 	return false, fmt.Errorf("bad response: %s", resp.Status)
	// }

	return false, nil
}

func (r *MovieRepository) GetCategoryID(ctx context.Context, categoryID int) (int, error) {
	// body := map[string]interface{}{
	// 	"name_tm":  "test",
	// 	"name_ru":  "test",
	// 	"name_en":  "test",
	// 	"belet_id": categoryID,
	// }

	// bodyBytes, err := json.Marshal(body)

	// if err != nil {
	// 	return 0, fmt.Errorf("failed to marshal request body: %w", err)
	// }

	// req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf(getCategoryIDURL, categoryID), bytes.NewBuffer(bodyBytes))

	// if err != nil {
	// 	return 0, fmt.Errorf("failed to create request: %w", err)
	// }

	// req.Header.Set("Authorization", r.accessToken)

	// client := &http.Client{
	// 	Timeout: time.Second * 10,
	// }

	// resp, err := client.Do(req)

	// if err != nil {
	// 	return 0, fmt.Errorf("request failed: %w", err)
	// }

	// defer resp.Body.Close()

	// if resp.StatusCode != http.StatusOK {
	// 	return 0, fmt.Errorf("bad response: %s", resp.Status)
	// }

	// var category models.GetIDResponse
	// err = json.NewDecoder(resp.Body).Decode(&category)

	// if err != nil {
	// 	return 0, fmt.Errorf("failed to decode response: %w", err)
	// }

	// return category.ID, nil
	return 1, nil
}

func (r *MovieRepository) GetGenreIDs(ctx context.Context, genres []string) ([]int, error) {
	genreIDs := []int{1, 2}

	for i := range genres {
		body := map[string]interface{}{
			"name_tm": genres[i],
			"name_ru": genres[i],
			"name_en": genres[i],
		}

		bodyBytes, err := json.Marshal(body)

		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, getGenreIDURL, bytes.NewBuffer(bodyBytes))

		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Authorization", r.accessToken)

		client := &http.Client{
			Timeout: time.Second * 10,
		}

		resp, err := client.Do(req)

		if err != nil {
			return nil, fmt.Errorf("request failed: %w", err)
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("bad response: %s", resp.Status)
		}

		var genre models.GetIDResponse
		err = json.NewDecoder(resp.Body).Decode(&genre)

		if err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}

		genreIDs = append(genreIDs, genre.ID)
	}

	return genreIDs, nil
}

func (r *MovieRepository) GetCountryIDs(ctx context.Context, countries []models.Country) ([]int, error) {
	countryIDs := []int{1, 2}
	for i := range countries {
		body := map[string]interface{}{
			"name_tm": countries[i].Name,
			"name_ru": countries[i].Name,
			"name_en": countries[i].Name,
		}

		bodyBytes, err := json.Marshal(body)

		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, getCountryIDURL, bytes.NewBuffer(bodyBytes))

		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Authorization", r.accessToken)

		client := &http.Client{
			Timeout: time.Second * 10,
		}

		resp, err := client.Do(req)

		if err != nil {
			return nil, fmt.Errorf("request failed: %w", err)
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("bad response: %s", resp.Status)
		}

		var country models.GetIDResponse
		err = json.NewDecoder(resp.Body).Decode(&country)

		if err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}

		countryIDs = append(countryIDs, country.ID)
	}

	return countryIDs, nil
}

func (r *MovieRepository) GetActorIDs(ctx context.Context, actors []models.Person) ([]int, error) {
	actorIDs := []int{1, 2}

	for i := range actors {
		body := map[string]interface{}{
			"name_tm": actors[i].Name,
			"name_ru": actors[i].Name,
			"name_en": actors[i].Name,
		}

		bodyBytes, err := json.Marshal(body)

		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, getActorIDURL, bytes.NewBuffer(bodyBytes))

		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Authorization", r.accessToken)

		client := &http.Client{
			Timeout: time.Second * 10,
		}

		resp, err := client.Do(req)

		if err != nil {
			return nil, fmt.Errorf("request failed: %w", err)
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("bad response: %s", resp.Status)
		}

		var actor models.GetIDResponse
		err = json.NewDecoder(resp.Body).Decode(&actor)

		if err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}

		actorIDs = append(actorIDs, actor.ID)
	}

	return actorIDs, nil
}
