package domain

import "github.com/docker/distribution/uuid"

type Session struct {
	ID     uuid.UUID
	UserID uuid.UUID

	OBUData []OBUData
}
