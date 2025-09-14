package app

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/rauan06/realtime-map/analytics/internal/usecase/analytics"
	"github.com/rauan06/realtime-map/go-commons/pkg/kafka/consumer"
	"github.com/rauan06/realtime-map/go-commons/pkg/logger"
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

	uc := analytics.New()

	kc, err := consumer.New(c, uc, *l, "myTopic")
	if err != nil {
		l.Fatal(err)
	}
	
	err = kc.Run()
	if err != nil {
		l.Fatal(err)
	}
}
