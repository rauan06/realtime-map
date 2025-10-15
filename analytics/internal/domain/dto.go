package domain

import (
	"time"
)

type (
	ObuDTO struct {
		SessionID string    `json:"session_id"`
		Long      float64   `json:"long"`
		Lat       float64   `json:"lat"`
		CreatedAt time.Time `json:"created_at"`
	}

	SessionDTO struct {
		SessionID string `json:"session_id"`
	}
)
