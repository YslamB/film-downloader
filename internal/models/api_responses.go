package models

type GetIDResponse struct {
	ID int `json:"id"`
}

// SearchResult represents the structure of searchResult.json
type SearchResult struct {
	Films   []SearchFilm `json:"films"`
	Message string       `json:"message"`
	Status  string       `json:"status"`
}

type SearchFilm struct {
	ID     int `json:"id"`
	TypeID int `json:"type_id"`
}
