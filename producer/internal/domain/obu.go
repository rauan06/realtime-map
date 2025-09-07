package domain

import (
	"time"
)

type OBUData struct {
	// ID        uuid.UUID
	Long      float64
	Lat       float64
	CreatedAt time.Time
}
