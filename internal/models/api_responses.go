package models

type GetIDResponse struct {
	ID int `json:"id"`
}

type SearchResult struct {
	Films   []SearchFilm `json:"films"`
	Message string       `json:"message"`
	Status  string       `json:"status"`
}

type SearchFilm struct {
	ID     int    `json:"id"`
	TypeID int    `json:"type_id"`
	Name   string `json:"name"`
}
