package models

type SearchRequest struct {
	Page  int    `json:"page"`
	Order string `json:"order"`
	Sort  string `json:"sort"`
}
