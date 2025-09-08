package domain

import (
	"time"
)

type OBUData struct {
	ID        []byte
	Long      float64
	Lat       float64
	Timestamp time.Time
}
