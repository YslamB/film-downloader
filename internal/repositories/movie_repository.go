package repositories

import (
	"bytes"
	"context"
	"encoding/json"
	"film-downloader/internal/config"
	"film-downloader/internal/models"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type MovieRepository struct {
	cfg *config.Config
}

const (
	refreshTokenURL       = "https://api.belet.tm/api/v1/auth/refresh"
	getMovieURL           = "http://95.85.126.217:5050/api/v1/admin/movies/ext/%s"
	createMovieURL        = "http://95.85.126.217:5050/api/v1/admin/movies"
	createSeasonURL       = "http://95.85.126.217:5050/api/v1/admin/seasons"
	createEpisodeURL      = "http://95.85.126.217:5050/api/v1/admin/seasons/episodes"
	getCategoryIDURL      = "http://95.85.126.217:5050/api/v1/admin/catalogs/categories"
	getGenreIDURL         = "http://95.85.126.217:5050/api/v1/admin/catalogs/genres"
	getCountryIDURL       = "http://95.85.126.217:5050/api/v1/admin/catalogs/countries"
	getActorIDURL         = "http://95.85.126.217:5050/api/v1/admin/catalogs/persons"
	getStudioIDURL        = "http://95.85.126.217:5050/api/v1/admin/catalogs/studios"
	getLanguageIDURL      = "http://95.85.126.217:5050/api/v1/admin/catalogs/languages"
	sendImageURL          = "http://95.85.126.217:5050/api/v1/admin/movies/images"
	updateActorURL        = "http://95.85.126.217:5050/api/v1/admin/catalogs/persons/%d"
	createMovieFileURL    = "http://95.85.126.217:5050/api/v1/admin/movies/file"
	assignMovieFileURL    = "http://95.85.126.217:5050/api/v1/admin/movies/files"
	checkEpisodeExistsURL = "http://95.85.126.217:5050/api/v1/admin/seasons/episodes/ext/%d"
)

func NewMovieRepository(cfg *config.Config) *MovieRepository {
	return &MovieRepository{
		cfg: cfg,
	}
}

// Helper function to decode response that could be either a number or an object with ID
func decodeIDResponse(bodyBytes []byte) (int, error) {
	// Try to unmarshal as GetIDResponse first
	var response models.GetIDResponse
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		// If that fails, try to unmarshal as a plain number
		var id int
		if numErr := json.Unmarshal(bodyBytes, &id); numErr != nil {
			return 0, fmt.Errorf("failed to decode response as object or number: object_err=%w, number_err=%w, body=%s", err, numErr, string(bodyBytes))
		}
		return id, nil
	}
	return response.ID, nil
}

func (r *MovieRepository) CreateSeason(ctx context.Context, season models.Season, movieID string) (int, error) {
	fmt.Println("🎬 Creating season:", season.Name, "for movie ID:", movieID)

	// Convert movieID from string to int
	movieIDInt, err := strconv.Atoi(movieID)
	if err != nil {
		return 0, fmt.Errorf("invalid movieID format: %w", err)
	}

	body := map[string]any{
		"movie_id": movieIDInt,
		"number":   1,
		"title":    season.Name,
		"ext_id":   season.ID,
	}
	bodyBytes, err := json.Marshal(body)

	if err != nil {
		return 0, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, createSeasonURL, bytes.NewBuffer(bodyBytes))

	if err != nil {
		fmt.Println("sd8fiui")
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept", "application/json")

	client := &http.Client{
		Timeout: time.Second * 30,
	}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("sd90f8io")
		return 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("s89dji")
		return 0, fmt.Errorf("bad response: %s, body: %s", resp.Status, string(body))
	}

	bodyBytes, err = io.ReadAll(resp.Body)

	if err != nil {
		return 0, fmt.Errorf("failed to read response body: %w", err)
	}

	seasonID, err := decodeIDResponse(bodyBytes)

	if err != nil {
		fmt.Println("sdf89 - CreateSeason decode error")
		return 0, fmt.Errorf("failed to decode CreateSeason response: %w", err)
	}

	fmt.Println("✅ Season created successfully with ID:", seasonID)
	return seasonID, nil
}

func (r *MovieRepository) RefreshToken(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, refreshTokenURL, nil)

	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Cookie", r.cfg.Cookie)
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := client.Do(req)

	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad response: %s", resp.Status)
	}
	cookies := resp.Header["Set-Cookie"]

	if len(cookies) > 0 {
		fmt.Println("🍪 Received Set-Cookie headers:", cookies)
		cookie := ""

		for i := range cookies {
			cookie += cookies[i] + ";"
		}

		r.cfg.SetCookie(cookie)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var response models.RefreshTokenResponse

	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	r.cfg.SetAccessToken(response.AccessToken)
	return nil
}

func (r *MovieRepository) CheckMovieExists(ctx context.Context, movieID string) (int, error) {
	fmt.Println("🔍 Checking movie with ID:", movieID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(getMovieURL, movieID), nil)

	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)

	if err != nil {
		return 0, fmt.Errorf("request failed: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var response models.GetIDResponse
		body, err := io.ReadAll(resp.Body)

		if err != nil {
			return 0, fmt.Errorf("failed to read response body: %w", err)
		}

		response.ID, err = strconv.Atoi(string(body))
		if err != nil {
			return 0, fmt.Errorf("failed to convert response body to int: %w", err)
		}
		return response.ID, nil
	}

	return 0, nil
}

func (r *MovieRepository) GetCategoryID(ctx context.Context, categoryID int) (int, error) {
	fmt.Println("🔍 Getting category ID:", categoryID)
	body := map[string]interface{}{
		"name_tm": uuid.New().String(),
		"name_ru": uuid.New().String(),
		"name_en": uuid.New().String(),
		"ext_id":  categoryID,
	}

	bodyBytes, err := json.Marshal(body)

	if err != nil {
		return 0, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, getCategoryIDURL, bytes.NewBuffer(bodyBytes))

	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", r.cfg.GetAccessToken())

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := client.Do(req)

	if err != nil {
		return 0, fmt.Errorf("request failed: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("bad response: %s", resp.Status)
	}

	var category models.GetIDResponse
	err = json.NewDecoder(resp.Body).Decode(&category)

	if err != nil {
		return 0, fmt.Errorf("failed to decode response s98d8jn: %w", err)
	}

	return category.ID, nil
}

func (r *MovieRepository) GetGenreIDs(ctx context.Context, genres []string) ([]int, error) {
	fmt.Println("🔍 Getting genre IDs:", genres)
	genreIDs := []int{}

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

		req.Header.Set("Authorization", r.cfg.GetAccessToken())

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
			return nil, fmt.Errorf("failed to decode response s8d9fu98: %w", err)
		}

		genreIDs = append(genreIDs, genre.ID)
	}

	return genreIDs, nil
}

func (r *MovieRepository) GetCountryIDs(ctx context.Context, countries []models.Country) ([]int, error) {
	fmt.Println("🔍 Getting country IDs:", countries)
	countryIDs := []int{}
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

		req.Header.Set("Authorization", r.cfg.GetAccessToken())

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
			return nil, fmt.Errorf("failed to decode response s98dfuj: %w", err)
		}

		countryIDs = append(countryIDs, country.ID)
	}

	return countryIDs, nil
}

func (r *MovieRepository) SendActorImage(ctx context.Context, actor models.Person) (int, error) {

	var imageInfo struct {
		URL    string `json:"url"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	}

	if err := json.Unmarshal([]byte(actor.Image), &imageInfo); err != nil {
		return 0, fmt.Errorf("failed to parse image JSON: %w", err)
	}

	imageData, filename, err := r.downloadImage(ctx, imageInfo.URL)
	if err != nil {
		return 0, fmt.Errorf("failed to download image: %w", err)
	}
	defer imageData.Close()

	imageID, err := r.uploadImage(ctx, imageData, filename, imageInfo.Width, imageInfo.Height)
	if err != nil {
		return 0, fmt.Errorf("failed to upload image: %w", err)
	}

	return imageID, nil
}

func (r *MovieRepository) SendMovieImage(ctx context.Context, image models.ImageSize) (int, error) {
	imageData, filename, err := r.downloadImage(ctx, image.URL)
	if err != nil {
		return 0, fmt.Errorf("failed to download image: %w", err)
	}
	defer imageData.Close()

	imageID, err := r.uploadImage(ctx, imageData, filename, image.Width, image.Height)
	if err != nil {
		return 0, fmt.Errorf("failed to upload image: %w", err)
	}

	return imageID, nil
}

func (r *MovieRepository) GetActorIDs(ctx context.Context, actors []models.Person) ([]int, error) {
	fmt.Println("🔍 Getting actor IDs:", actors)
	actorIDs := []int{}

	for i := range actors {
		body := map[string]interface{}{
			"bio":       actors[i].Name,
			"full_name": actors[i].Name,
			"image_id":  1,
		}

		bodyBytes, err := json.Marshal(body)

		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, getActorIDURL, bytes.NewBuffer(bodyBytes))

		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Authorization", r.cfg.GetAccessToken())

		client := &http.Client{
			Timeout: time.Second * 10,
		}

		resp, err := client.Do(req)

		if err != nil {
			return nil, fmt.Errorf("request failed: %w", err)
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusConflict {
			return nil, fmt.Errorf("bad response: %s", resp.Status)
		}

		var actorRes models.GetIDResponse
		err = json.NewDecoder(resp.Body).Decode(&actorRes)

		if err != nil {
			return nil, fmt.Errorf("failed to decode response soidfj: %w", err)
		}

		if resp.StatusCode == http.StatusConflict {
			actorIDs = append(actorIDs, actorRes.ID)
		}

		if resp.StatusCode == http.StatusOK {
			imageID, err := r.SendActorImage(ctx, actors[i])

			if err != nil {
				return nil, fmt.Errorf("failed to send actor image: %w", err)
			}

			err = r.UpdateActorImage(ctx, actorRes.ID, actors[i], imageID)

			if err != nil {
				return nil, fmt.Errorf("failed to update actor image: %w", err)
			}

			actorIDs = append(actorIDs, actorRes.ID)
		}
	}

	return actorIDs, nil
}

func (r *MovieRepository) UpdateActorImage(ctx context.Context, actorID int, actor models.Person, imageID int) error {
	fmt.Println("🔍 Updating actor image:", actorID, actor.Name, imageID)

	body := map[string]interface{}{
		"bio":       actor.Name,
		"full_name": actor.Name,
		"image_id":  imageID,
	}

	bodyBytes, err := json.Marshal(body)

	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, fmt.Sprintf(updateActorURL, actorID), bytes.NewBuffer(bodyBytes))

	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", r.cfg.GetAccessToken())
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept", "application/json")

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := client.Do(req)

	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("bad response: %s, body: %s", resp.Status, string(body))
	}

	return nil
}

func (r *MovieRepository) downloadImage(ctx context.Context, imageURL string) (io.ReadCloser, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, imageURL, nil)

	if err != nil {
		return nil, "", fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{
		Timeout: time.Second * 30,
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, "", fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, "", fmt.Errorf("bad response: %s", resp.Status)
	}

	filename := filepath.Base(imageURL)

	if filename == "." || filename == "/" {
		filename = "actor_image.jpg"
	}

	if !strings.Contains(filename, ".") {
		filename += ".jpg"
	}

	return resp.Body, filename, nil
}

func (r *MovieRepository) uploadImage(ctx context.Context, imageData io.ReadCloser, filename string, width, height int) (int, error) {

	imageBytes, err := io.ReadAll(imageData)
	if err != nil {
		return 0, fmt.Errorf("failed to read image data: %w", err)
	}

	contentType := http.DetectContentType(imageBytes)

	if !strings.HasPrefix(contentType, "image/") {
		return 0, fmt.Errorf("invalid file type: %s, expected image", contentType)
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	h := make(map[string][]string)
	h["Content-Disposition"] = []string{fmt.Sprintf(`form-data; name="image"; filename="%s"`, filename)}
	h["Content-Type"] = []string{contentType}

	fileWriter, err := writer.CreatePart(h)
	if err != nil {
		return 0, fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err := fileWriter.Write(imageBytes); err != nil {
		return 0, fmt.Errorf("failed to write image data: %w", err)
	}

	if err := writer.WriteField("width", fmt.Sprintf("%d", width)); err != nil {
		return 0, fmt.Errorf("failed to write width field: %w", err)
	}

	if err := writer.WriteField("height", fmt.Sprintf("%d", height)); err != nil {
		return 0, fmt.Errorf("failed to write height field: %w", err)
	}

	if err := writer.Close(); err != nil {
		return 0, fmt.Errorf("failed to close writer: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, sendImageURL, &buf)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("accept", "application/json")

	client := &http.Client{
		Timeout: time.Second * 30,
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("bad response: %s, body: %s", resp.Status, string(body))
	}

	var response struct {
		ID int `json:"id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, fmt.Errorf("failed to decode response sd8f9: %w", err)
	}

	return response.ID, nil
}

func (r *MovieRepository) GetStudioIDs(ctx context.Context, studios []models.Studio) ([]int, error) {
	fmt.Println("🔍 Getting studio IDs:", studios)
	studioIDs := []int{}

	for i := range studios {
		body := map[string]interface{}{
			"name": studios[i].Name,
		}

		bodyBytes, err := json.Marshal(body)

		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, getStudioIDURL, bytes.NewBuffer(bodyBytes))

		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("accept", "application/json")

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

		var studio models.GetIDResponse
		err = json.NewDecoder(resp.Body).Decode(&studio)

		if err != nil {
			return nil, fmt.Errorf("failed to decode response sd9f8ji: %w", err)
		}

		studioIDs = append(studioIDs, studio.ID)
	}

	return studioIDs, nil
}

func (r *MovieRepository) GetLanguageID(ctx context.Context, language string) (int, error) {
	fmt.Println("🔍 Getting language ID for:", language)

	body := map[string]interface{}{
		"name_tm": language,
		"name_ru": language,
		"name_en": language,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, getLanguageIDURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", r.cfg.GetAccessToken())
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept", "application/json")

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("bad response: %s", resp.Status)
	}

	var language_response models.GetIDResponse
	err = json.NewDecoder(resp.Body).Decode(&language_response)
	if err != nil {
		return 0, fmt.Errorf("failed to decode response sduf8h98s: %w", err)
	}

	return language_response.ID, nil
}

func (r *MovieRepository) CreateMovie(ctx context.Context, movie models.Film, genreIDs, countryIDs, actorIDs, directorIDs, studioIDs []int, languageID, verticalImageID, verticalWithoutNameImageID, horizontalWithNameImageID, horizontalWithoutNameImageID, nameImageID int) (int, error) {
	fmt.Println("🎬 Creating movie:", movie.Name)

	duration := 0

	if movie.Duration != "" {
		duration, _ = strconv.Atoi(movie.Duration)
	} else {
		duration = 90
	}

	ageRestriction := 0

	if movie.Age != "" {
		ageRestriction, _ = strconv.Atoi(movie.Age)
	} else {
		ageRestriction = 16
	}

	movieType := "movie"

	if movie.TypeID == 2 {
		movieType = "series"
	}

	body := map[string]any{
		"ext_id":                  movie.ID,
		"title":                   movie.Name,
		"description":             movie.Description,
		"release_year":            movie.Year,
		"duration":                duration,
		"age_restriction":         ageRestriction,
		"rating":                  movie.RatingKP,
		"rating_imdb":             movie.RatingIMDB,
		"rating_kinopoisk":        movie.RatingKP,
		"category_id":             movie.CategoryID,
		"type":                    movieType,
		"language_id":             languageID,
		"color":                   "",
		"genre_ids":               genreIDs,
		"country_ids":             countryIDs,
		"actor_ids":               actorIDs,
		"director_ids":            directorIDs,
		"studio_ids":              studioIDs,
		"vertical":                verticalImageID,
		"vertical_without_name":   verticalWithoutNameImageID,
		"horizontal_with_name":    horizontalWithNameImageID,
		"horizontal_without_name": horizontalWithoutNameImageID,
		"image_name":              nameImageID,
	}

	bodyBytes, err := json.Marshal(body)

	if err != nil {
		return 0, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, createMovieURL, bytes.NewBuffer(bodyBytes))

	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", r.cfg.GetAccessToken())
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept", "application/json")

	client := &http.Client{
		Timeout: time.Second * 30,
	}

	resp, err := client.Do(req)

	if err != nil {
		return 0, fmt.Errorf("request failed: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("bad response: %s, body: %s", resp.Status, string(body))
	}

	unmarshalBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return 0, fmt.Errorf("failed to read response body: %w", err)
	}

	movieID, err := decodeIDResponse(unmarshalBody)
	if err != nil {
		return 0, fmt.Errorf("failed to decode CreateMovie response: %w", err)
	}

	fmt.Println("✅ Movie created successfully with ID:", movieID)
	return movieID, nil
}

func (r *MovieRepository) GetFileID(ctx context.Context, name string) (int, error) {
	fmt.Println("🎬 Creating movie file for:", name)

	body := map[string]any{
		"path": "movies/" + name + "/master.m3u8",
		"type": "application/x-mpegURL",
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, createMovieFileURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", r.cfg.GetAccessToken())
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept", "application/json")

	client := &http.Client{
		Timeout: time.Second * 30,
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("bad response: %s, body: %s", resp.Status, string(body))
	}

	var response models.GetIDResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, fmt.Errorf("failed to decode response s98dfu98sd: %w", err)
	}

	fmt.Println("✅ Movie file created successfully with ID:", response.ID)
	return response.ID, nil
}

func (r *MovieRepository) CreateMovieFile(ctx context.Context, fileID, movieID int) error {
	fmt.Println("🔗 Assigning file ID", fileID, "to movie ID", movieID)

	body := map[string]any{
		"file_id":  fileID,
		"movie_id": movieID,
	}

	bodyBytes, err := json.Marshal(body)

	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, assignMovieFileURL, bytes.NewBuffer(bodyBytes))

	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept", "application/json")

	client := &http.Client{
		Timeout: time.Second * 30,
	}

	resp, err := client.Do(req)

	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("bad response: %s, body: %s", resp.Status, string(body))
	}

	fmt.Println("✅ Successfully assigned file to movie")
	return nil
}

func (r *MovieRepository) CreateEpisode(ctx context.Context, episode models.Episode, seasonID int) (int, error) {
	fmt.Println("📺 Creating episode:", episode.Name, "for season ID:", seasonID)

	var imageInfo struct {
		URL    string `json:"url"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	}
	imageInfo.URL = episode.Image
	imageInfo.Width = 640
	imageInfo.Height = 360
	imageData, filename, err := r.downloadImage(ctx, imageInfo.URL)

	if err != nil {
		return 0, fmt.Errorf("failed to download image: %w", err)
	}
	defer imageData.Close()
	imageID, err := r.uploadImage(ctx, imageData, filename, imageInfo.Width, imageInfo.Height)

	if err != nil {
		return 0, fmt.Errorf("failed to upload image: %w", err)
	}
	duration := 30

	if episode.Duration != "" {
		if parsed, err := strconv.Atoi(episode.Duration); err == nil {
			duration = parsed
		}
	}

	body := map[string]any{
		"duration":  duration,
		"file_id":   episode.FileID,
		"image_id":  imageID,
		"number":    1,
		"ext_id":    episode.ID,
		"season_id": seasonID,
		"title":     episode.Name,
	}
	bodyBytes, err := json.Marshal(body)

	if err != nil {
		return 0, fmt.Errorf("failed to marshal request body: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, createEpisodeURL, bytes.NewBuffer(bodyBytes))

	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", r.cfg.GetAccessToken())
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept", "application/json")

	client := &http.Client{
		Timeout: time.Second * 30,
	}
	resp, err := client.Do(req)

	if err != nil {
		return 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("bad response: %s, body: %s", resp.Status, string(body))
	}

	if resp.StatusCode == http.StatusConflict {
		return 0, fmt.Errorf("episode already exists")
	}
	var response models.GetIDResponse

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, fmt.Errorf("failed to decode response sd89fus9i8df: %w", err)
	}
	fmt.Println("✅ Episode created successfully with ID:", response.ID)
	return response.ID, nil
}

func (r *MovieRepository) CheckEpisodeExists(ctx context.Context, extID int) error {
	url := fmt.Sprintf(checkEpisodeExistsURL, extID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)

	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return fmt.Errorf("episode with external ID %d already exists", extID)
	}
	return nil
}
