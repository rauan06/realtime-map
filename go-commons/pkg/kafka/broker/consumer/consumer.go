package consumer

import (
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/goccy/go-json"

	"github.com/evrone/go-clean-template/pkg/logger"
)

const (
	_defaultWaitTime = 5 * time.Second
	_defaultAttempts = 10
	_defaultTimeout  = 2 * time.Second
)

// CallHandler -.
type CallHandler func(*kafka.Message) (interface{}, error)

// Consumer -.
type Consumer struct {
	consumer *kafka.Consumer
	error    chan error
	stop     chan struct{}
	router   map[string]CallHandler

	brokers  string
	groupID  string
	topic    string
	timeout  time.Duration
	attempts int
	waitTime time.Duration

	logger logger.Interface
}

// New -.
func New(brokers, groupID, topic string, router map[string]CallHandler, l logger.Interface, opts ...Option) (*Consumer, error) {
	c := &Consumer{
		error:    make(chan error),
		stop:     make(chan struct{}),
		router:   router,
		brokers:  brokers,
		groupID:  groupID,
		topic:    topic,
		timeout:  _defaultTimeout,
		attempts: _defaultAttempts,
		waitTime: _defaultWaitTime,
		logger:   l,
	}

	// Custom options
	for _, opt := range opts {
		opt(c)
	}

	err := c.AttemptConnect()
	if err != nil {
		return nil, fmt.Errorf("kafka consumer - New - c.AttemptConnect: %w", err)
	}

	return c, nil
}

// AttemptConnect -.
func (c *Consumer) AttemptConnect() error {
	var err error
	for i := c.attempts; i > 0; i-- {
		if err = c.connect(); err == nil {
			break
		}

		c.logger.Warn(fmt.Sprintf("Kafka consumer is trying to connect, attempts left: %d, error: %v", i, err))
		time.Sleep(c.waitTime)
	}

	if err != nil {
		return fmt.Errorf("kafka consumer - AttemptConnect - c.connect: %w", err)
	}

	return nil
}

func (c *Consumer) connect() error {
	var err error

	c.consumer, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": c.brokers,
		"group.id":          c.groupID,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return fmt.Errorf("kafka.NewConsumer: %w", err)
	}

	err = c.consumer.SubscribeTopics([]string{c.topic}, nil)
	if err != nil {
		return fmt.Errorf("c.consumer.SubscribeTopics: %w", err)
	}

	return nil
}

// Start -.
func (c *Consumer) Start() {
	go c.consume()
}

func (c *Consumer) consume() {
	for {
		select {
		case <-c.stop:
			return
		default:
			msg, err := c.consumer.ReadMessage(-1)
			if err != nil {
				c.logger.Error(err, "kafka consumer - consume - ReadMessage")
				c.reconnect()
				return
			}

			c.serveCall(msg)
		}
	}
}

func (c *Consumer) serveCall(msg *kafka.Message) {
	var msgType string
	for _, h := range msg.Headers {
		if h.Key == "type" {
			msgType = string(h.Value)
			break
		}
	}

	callHandler, ok := c.router[msgType]
	if !ok {
		c.logger.Warn(fmt.Sprintf("No handler found for message type: %s", msgType))
		return
	}

	response, err := callHandler(msg)
	if err != nil {
		c.logger.Error(err, "kafka consumer - serveCall - handler error")
		return
	}

	// optional: encode response (e.g. for forwarding)
	body, err := json.Marshal(response)
	if err != nil {
		c.logger.Error(err, "kafka consumer - serveCall - json.Marshal")
		return
	}

	c.logger.Info(fmt.Sprintf("Processed message type=%s key=%s response=%s",
		msgType, string(msg.Key), string(body)))
}

func (c *Consumer) reconnect() {
	close(c.stop)

	err := c.AttemptConnect()
	if err != nil {
		c.error <- err
		close(c.error)
		return
	}

	c.stop = make(chan struct{})
	go c.consume()
}

// Notify -.
func (c *Consumer) Notify() <-chan error {
	return c.error
}

// Shutdown -.
func (c *Consumer) Shutdown() error {
	select {
	case <-c.error:
		return nil
	default:
	}

	close(c.stop)
	time.Sleep(c.timeout)

	err := c.consumer.Close()
	if err != nil {
		return fmt.Errorf("kafka consumer - Shutdown - c.consumer.Close: %w", err)
	}

	return nil
}
