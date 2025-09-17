package domain

import (
	"time"
)

type (
	OBU_dto struct {
		SessionID string    `json:"session_id"`
		Long      float64   `json:"long"`
		Lat       float64   `json:"lat"`
		CreatedAt time.Time `json:"created_at"`
	}

	Session_dto struct {
		SessionID string `json:"session_id"`
	}
)
