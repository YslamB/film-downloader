package models

import "time"

var (
	Resolutions = map[string]string{
		"1080p": "1920x1080",
		"720p":  "1280x720",
		"480p":  "854x480",
	}

	Bandwidths = map[string]string{
		"1080p": "5128000",
		"720p":  "1500000",
		"480p":  "5128000",
	}

	Codecs = map[string]string{
		"1080p": "avc1.640028,mp4a.40.2",
		"720p":  "avc1.42001e,mp4a.40.2",
		"480p":  "avc1.42001e,mp4a.40.2",
	}
)

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
