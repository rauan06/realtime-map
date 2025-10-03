package app

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/rauan06/realtime-map/analytics/internal/domain"
	repo_postgres "github.com/rauan06/realtime-map/analytics/internal/repo/postgres"
	"github.com/rauan06/realtime-map/analytics/internal/usecase/analytics"
	"github.com/rauan06/realtime-map/go-commons/pkg/kafka/consumer"
	"github.com/rauan06/realtime-map/go-commons/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Run() {
	l := logger.New("DEBUG")

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		l.Fatal(err)
	}

	dsn := "host=db user=postgres password=example dbname=realtimedb port=5432 sslmode=disable TimeZone=Asia/Almaty"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	obu_repo := repo_postgres.New[domain.OBUData](db)
	session_repo := repo_postgres.New[domain.Session](db)

	uc := analytics.New(*l, obu_repo, session_repo)

	kc, err := consumer.New(c, uc, *l, "myTopic")
	if err != nil {
		l.Fatal(err)
	}

	err = kc.Run()
	if err != nil {
		l.Fatal(err)
	}
}
