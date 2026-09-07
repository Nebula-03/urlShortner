package model

import "time"

type URL struct {
	ID          int64
	Alias       string
	OriginalURL string
	CreatedAt   time.Time
}
