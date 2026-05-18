package usecase

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"github.com/rauan06/realtime-map/go-commons/pkg/logger"
	"github.com/rauan06/realtime-map/notification/internal/geofence"
	"github.com/rauan06/realtime-map/notification/internal/notifier"
)

// Geofencing is the kafka-driven enter/exit detector. State is per (layer,
// source_id) → set of fence names the entity is currently inside, so we emit
// at-most-one alert per transition.
type Geofencing struct {
	registry *geofence.Registry
	notifier *notifier.Notifier
	l        logger.Logger

	mu    sync.Mutex
	state map[string]map[string]struct{}
}

func New(reg *geofence.Registry, n *notifier.Notifier, l logger.Logger) *Geofencing {
	return &Geofencing{
		registry: reg,
		notifier: n,
		l:        l,
		state:    make(map[string]map[string]struct{}),
	}
}

// ProcessMessage matches the go-commons consumer's `uc` interface so the
// service can plug directly into the existing KafkaConsumer.
func (uc *Geofencing) ProcessMessage(msg *kafka.Message) {
	if msg == nil || msg.TopicPartition.Topic == nil {
		return
	}
	layer := topicToLayer(*msg.TopicPartition.Topic)

	var payload map[string]interface{}
	if err := json.Unmarshal(msg.Value, &payload); err != nil {
		uc.l.Error("notification: unmarshal: %v", err)
		return
	}

	lat, lng, ok := extractLatLng(payload)
	if !ok {
		return
	}

	sourceID := extractSourceID(payload, msg.Key)
	if sourceID == "" {
		return
	}

	matches := uc.registry.Match(layer, lat, lng)
	key := layer + ":" + sourceID

	uc.mu.Lock()
	prev := uc.state[key]
	if prev == nil {
		prev = map[string]struct{}{}
	}
	cur := make(map[string]struct{}, len(matches))
	for _, m := range matches {
		cur[m] = struct{}{}
	}
	uc.state[key] = cur
	uc.mu.Unlock()

	now := time.Now().UTC()

	for name := range cur {
		if _, was := prev[name]; was {
			continue
		}
		uc.notifier.Dispatch(context.Background(), notifier.Alert{
			Event: "enter", Layer: layer, SourceID: sourceID, Fence: name,
			Lat: lat, Lng: lng, At: now,
		})
	}
	for name := range prev {
		if _, still := cur[name]; still {
			continue
		}
		uc.notifier.Dispatch(context.Background(), notifier.Alert{
			Event: "exit", Layer: layer, SourceID: sourceID, Fence: name,
			Lat: lat, Lng: lng, At: now,
		})
	}
}

func topicToLayer(topic string) string {
	switch topic {
	case "etl_flights":
		return "flight"
	case "etl_ships":
		return "ship"
	case "etl_transport":
		return "transport"
	case "obu_data":
		return "obu"
	default:
		return topic
	}
}

func extractLatLng(p map[string]interface{}) (float64, float64, bool) {
	lat, ok1 := toFloat(p["lat"])
	if !ok1 {
		lat, ok1 = toFloat(p["latitude"])
	}
	lng, ok2 := toFloat(p["lng"])
	if !ok2 {
		lng, ok2 = toFloat(p["long"])
	}
	if !ok2 {
		lng, ok2 = toFloat(p["longitude"])
	}
	return lat, lng, ok1 && ok2
}

func extractSourceID(p map[string]interface{}, kafkaKey []byte) string {
	candidates := []string{"source_id", "device_id", "session_id", "icao24", "mmsi", "vehicle_id", "road_id"}
	for _, k := range candidates {
		if v, ok := p[k]; ok {
			if s, ok := v.(string); ok && s != "" {
				return s
			}
		}
	}
	if len(kafkaKey) > 0 {
		return string(kafkaKey)
	}
	return ""
}

func toFloat(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	case json.Number:
		f, err := n.Float64()
		return f, err == nil
	default:
		return 0, false
	}
}
