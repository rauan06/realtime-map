package broker

import (
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// Config -.
type Config struct {
	Brokers  string
	GroupID  string
	WaitTime time.Duration
	Attempts int
}

// Connection -.
type Connection struct {
	Topic     string
	Config
	Consumer *kafka.Consumer
	Producer *kafka.Producer
	Events   chan *kafka.Message
}

// New -.
func New(topic string, cfg Config) *Connection {
	return &Connection{
		Topic:   topic,
		Config:  cfg,
		Events:  make(chan *kafka.Message),
	}
}

// AttemptConnect -.
func (c *Connection) AttemptConnect() error {
	var err error
	for i := c.Attempts; i > 0; i-- {
		if err = c.connect(); err == nil {
			break
		}

		log.Printf("Kafka is trying to connect, attempts left: %d, error: %v", i, err)
		time.Sleep(c.WaitTime)
	}

	if err != nil {
		return fmt.Errorf("kafka_rpc - AttemptConnect - c.connect: %w", err)
	}

	return nil
}

func (c *Connection) connect() error {
	var err error

	// Create Consumer
	c.Consumer, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": c.Brokers,
		"group.id":          c.GroupID,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return fmt.Errorf("kafka.NewConsumer: %w", err)
	}

	err = c.Consumer.Subscribe(c.Topic, nil)
	if err != nil {
		return fmt.Errorf("c.Consumer.Subscribe: %w", err)
	}

	// Create Producer
	c.Producer, err = kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": c.Brokers,
	})
	if err != nil {
		return fmt.Errorf("kafka.NewProducer: %w", err)
	}

	// Start listening to consumer messages
	go func() {
		for {
			msg, err := c.Consumer.ReadMessage(-1)
			if err == nil {
				c.Events <- msg
			} else {
				log.Printf("Consumer error: %v\n", err)
			}
		}
	}()

	return nil
}

// SendMessage -.
func (c *Connection) SendMessage(key, value []byte) error {
	return c.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &c.Topic, Partition: kafka.PartitionAny},
		Key:            key,
		Value:          value,
	}, nil)
}
