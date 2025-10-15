package analytics

import (
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"

	"github.com/rauan06/realtime-map/analytics/internal/domain"
	"github.com/rauan06/realtime-map/analytics/internal/repo"
	"github.com/rauan06/realtime-map/analytics/internal/usecase"
	"github.com/rauan06/realtime-map/go-commons/pkg/logger"
)

type UseCase struct {
	obuDB     repo.IDatabase[domain.OBUData]
	sessionDB repo.IDatabase[domain.Session]

	l logger.Logger
}

func New(l logger.Logger, oduDB repo.IDatabase[domain.OBUData], sessionDB repo.IDatabase[domain.Session]) usecase.IAnalyticsUseCase {
	return &UseCase{
		obuDB:     oduDB,
		sessionDB: sessionDB,
		l:         l,
	}
}

func (uc *UseCase) ProcessMessage(msg *kafka.Message) {
	uc.l.Info("received msg: %+v, key: %s, value: %s\n", msg, msg.Key, msg.Value)

	var dto domain.ObuDTO
	if err := json.Unmarshal(msg.Value, &dto); err != nil {
		uc.l.Error("error marshaling obu data: %v", err)

		return
	}

	sessionID, err := uuid.Parse(dto.SessionID)
	if err != nil {
		uc.l.Error("error parsing session uuid: %v", err)

		return
	}

	data := &domain.OBUData{
		SessionID: sessionID,
		Point:     *domain.NewGeoPoint(dto.Long, dto.Lat),
	}

	if err := uc.obuDB.Create(data); err != nil {
		uc.l.Error("error inserting obu data: %v", err)

		return
	}

	uc.l.Info("successfully processed msg: %+v", msg)
}
