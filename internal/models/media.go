package models

import "time"

type MediaType string

const (
	Movie MediaType = "movie"
	Drama MediaType = "drama"
)

type MediaEntry struct {
	ID          int
	Title       string
	Type        MediaType
	Rating      float64
	Comment     string
	DateWatched time.Time
	CreatedAt   time.Time
}
