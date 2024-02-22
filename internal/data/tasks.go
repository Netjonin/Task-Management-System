package data

import (
	"time"
)

type Task struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"-"`
	Status      string    `json:"status"`
	ExpiredAt   time.Time `json:"expired_at"`
	Expired     bool      `json:"expired"`
	Version     int32     `json:"version"`
}
