package flight

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/rauan06/realtime-map/etl/internal/domain"
)

const minStateFields = 11

const apiURL = "https://opensky-network.org/api/states/all"

type Extractor struct {
	client *http.Client
}

func New(timeout time.Duration) *Extractor {
	return &Extractor{
		client: &http.Client{Timeout: timeout},
	}
}

type response struct {
	Time   int64           `json:"time"`
	States [][]interface{} `json:"states"`
}

func (e *Extractor) Extract(ctx context.Context) ([]domain.RawRecord, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("flight extract: new request: %w", err)
	}

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("flight extract: do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("flight extract: %w: %d", domain.ErrUpstreamStatus, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("flight extract: read body: %w", err)
	}

	var data response
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("flight extract: unmarshal: %w", err)
	}

	records := make([]domain.RawRecord, 0, len(data.States))

	for _, s := range data.States {
		if len(s) < minStateFields {
			continue
		}

		icao24, ok := s[0].(string)
		if !ok || icao24 == "" {
			continue
		}

		lon, ok1 := toFloat64(s[5])
		lat, ok2 := toFloat64(s[6])

		if !ok1 || !ok2 {
			continue
		}

		callsign, _ := s[1].(string) //nolint:errcheck // optional field
		altitude, _ := toFloat64(s[7])
		onGround, _ := s[8].(bool) //nolint:errcheck // optional field
		velocity, _ := toFloat64(s[9])
		heading, _ := toFloat64(s[10])

		records = append(records, domain.RawRecord{
			SourceID:  icao24,
			Timestamp: time.Unix(data.Time, 0),
			Fields: map[string]any{
				"icao24":    icao24,
				"callsign":  strings.TrimSpace(callsign),
				"lat":       lat,
				"lng":       lon,
				"altitude":  altitude,
				"on_ground": onGround,
				"velocity":  velocity,
				"heading":   heading,
			},
		})
	}

	return records, nil
}

func toFloat64(v any) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case json.Number:
		f, err := n.Float64()

		return f, err == nil
	default:
		return 0, false
	}
}
