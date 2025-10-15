package app

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/rauan06/realtime-map/analytics/internal/domain"
	repoPostgres "github.com/rauan06/realtime-map/analytics/internal/repo/postgres"
	"github.com/rauan06/realtime-map/analytics/internal/usecase/analytics"
	"github.com/rauan06/realtime-map/go-commons/pkg/kafka/consumer"
	"github.com/rauan06/realtime-map/go-commons/pkg/logger"
)

func Run() {
	l := logger.New("DEBUG")

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "broker:29092",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		l.Fatal(err)
	}

	dsn := "host=db user=postgres password=example dbname=realtimedb port=5432 sslmode=disable TimeZone=Asia/Almaty"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		l.Fatal(err)
	}

	obuRepo := repoPostgres.New[domain.OBUData](db)
	sessionRepo := repoPostgres.New[domain.Session](db)

	uc := analytics.New(*l, obuRepo, sessionRepo)

	kc, err := consumer.New(c, uc, *l, "myTopic")
	if err != nil {
		l.Fatal(err)
	}

	err = kc.Run()
	if err != nil {
		l.Fatal(err)
	}
}
