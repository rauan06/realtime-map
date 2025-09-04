package producer

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
)

// ErrConnectionClosed -.
var ErrConnectionClosed = errors.New("kafka producer - Produce - Connection closed")

const (
	_defaultWaitTime = 5 * time.Second
	_defaultAttempts = 10
	_defaultTimeout  = 2 * time.Second
)

// Message -.
type Message struct {
	Key         []byte
	Value       []byte
	Headers     []kafka.Header
	Correlation string
}

type pendingCall struct {
	done chan struct{}
	err  error
}

// Producer -.
type Producer struct {
	producer *kafka.Producer
	topic    string
	error    chan error
	stop     chan struct{}

	rw    sync.RWMutex
	calls map[string]*pendingCall

	brokers  string
	timeout  time.Duration
	attempts int
	waitTime time.Duration
}

// New -.
func New(brokers, topic string, opts ...Option) (*Producer, error) {
	p := &Producer{
		brokers:  brokers,
		topic:    topic,
		error:    make(chan error),
		stop:     make(chan struct{}),
		calls:    make(map[string]*pendingCall),
		timeout:  _defaultTimeout,
		attempts: _defaultAttempts,
		waitTime: _defaultWaitTime,
	}

	// Custom options
	for _, opt := range opts {
		opt(p)
	}

	err := p.AttemptConnect()
	if err != nil {
		return nil, fmt.Errorf("kafka producer - New - p.AttemptConnect: %w", err)
	}

	go p.deliveryReports()

	return p, nil
}

// AttemptConnect -.
func (p *Producer) AttemptConnect() error {
	var err error
	for i := p.attempts; i > 0; i-- {
		if err = p.connect(); err == nil {
			break
		}

		fmt.Printf("Kafka producer is trying to connect, attempts left: %d, error: %v\n", i, err)
		time.Sleep(p.waitTime)
	}

	if err != nil {
		return fmt.Errorf("kafka producer - AttemptConnect - p.connect: %w", err)
	}

	return nil
}

func (p *Producer) connect() error {
	var err error

	p.producer, err = kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": p.brokers,
	})
	if err != nil {
		return fmt.Errorf("kafka.NewProducer: %w", err)
	}

	return nil
}

// Produce -.
func (p *Producer) Produce(handler string, request interface{}) error {
	select {
	case <-p.stop:
		time.Sleep(p.timeout)
		select {
		case <-p.stop:
			return ErrConnectionClosed
		default:
		}
	default:
	}

	corrID := uuid.New().String()
	var requestBody []byte
	var err error

	if request != nil {
		requestBody, err = json.Marshal(request)
		if err != nil {
			return fmt.Errorf("json.Marshal: %w", err)
		}
	}

	call := &pendingCall{done: make(chan struct{})}
	p.addCall(corrID, call)
	defer p.deleteCall(corrID)

	err = p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
		Key:            []byte(handler),
		Value:          requestBody,
		Headers: []kafka.Header{
			{Key: "correlation_id", Value: []byte(corrID)},
		},
	}, nil)
	if err != nil {
		return fmt.Errorf("p.producer.Produce: %w", err)
	}

	select {
	case <-time.After(p.timeout):
		return fmt.Errorf("produce timeout: %w", err)
	case <-call.done:
		return call.err
	}
}

func (p *Producer) deliveryReports() {
	for e := range p.producer.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			corrID := ""
			for _, h := range ev.Headers {
				if h.Key == "correlation_id" {
					corrID = string(h.Value)
				}
			}
			if corrID == "" {
				continue
			}

			p.rw.RLock()
			call, ok := p.calls[corrID]
			p.rw.RUnlock()
			if !ok {
				continue
			}

			if ev.TopicPartition.Error != nil {
				call.err = ev.TopicPartition.Error
			}

			close(call.done)
		}
	}
}

func (p *Producer) addCall(corrID string, call *pendingCall) {
	p.rw.Lock()
	p.calls[corrID] = call
	p.rw.Unlock()
}

func (p *Producer) deleteCall(corrID string) {
	p.rw.Lock()
	delete(p.calls, corrID)
	p.rw.Unlock()
}

// Notify -.
func (p *Producer) Notify() <-chan error {
	return p.error
}

// Shutdown -.
func (p *Producer) Shutdown() error {
	select {
	case <-p.error:
		return nil
	default:
	}

	close(p.stop)
	time.Sleep(p.timeout)

	p.producer.Close()

	return nil
}
