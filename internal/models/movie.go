package models

import "time"

// BasicMovie represents a simplified movie structure for internal use
type BasicMovie struct {
	ID          string
	Title       string
	Description string
	ReleaseDate time.Time
	Rating      float64
	Genre       string
	Director    string
	Actor       string
}
