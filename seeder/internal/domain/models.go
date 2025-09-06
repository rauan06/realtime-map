package domain

import (
	"github.com/google/uuid"
)

type OBUData struct {
	ID   uuid.UUID
	Long float64
	Lat  float64
}
