package domain

import (
	"time"

	"github.com/google/uuid"
)

type OBUData struct {
	SessionID uuid.UUID `json:"session_id"`
	Long      float64   `json:"long"`
	Lat       float64   `json:"lat"`
	CreatedAt time.Time `json:"created_at"`
}
