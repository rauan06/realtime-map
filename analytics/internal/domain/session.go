package domain

import "github.com/google/uuid"

type Session struct {
	ID     uuid.UUID
	UserID uuid.UUID

	OBUData []OBUData
}

func (Session) TableName() string { return "sessions" }
