package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rauan06/realtime-map/etl/internal/domain"
)

const apiURL = "https://api-v3.mbta.com/vehicles"

type Extractor struct {
	client *http.Client
}

func New(timeout time.Duration) *Extractor {
	return &Extractor{
		client: &http.Client{Timeout: timeout},
	}
}

type response struct {
	Data []vehicle `json:"data"`
}

type vehicle struct {
	ID         string     `json:"id"`
	Attributes attributes `json:"attributes"`
}

type attributes struct {
	Latitude      *float64 `json:"latitude"`
	Longitude     *float64 `json:"longitude"`
	Bearing       *float64 `json:"bearing"`
	Speed         *float64 `json:"speed"`
	Label         string   `json:"label"`
	CurrentStatus string   `json:"current_status"`
	UpdatedAt     string   `json:"updated_at"`
}

func (e *Extractor) Extract(ctx context.Context) ([]domain.RawRecord, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("transport extract: new request: %w", err)
	}

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("transport extract: do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("transport extract: %w: %d", domain.ErrUpstreamStatus, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("transport extract: read body: %w", err)
	}

	var data response
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("transport extract: unmarshal: %w", err)
	}

	records := make([]domain.RawRecord, 0, len(data.Data))
	for _, v := range data.Data {
		if v.Attributes.Latitude == nil || v.Attributes.Longitude == nil {
			continue
		}

		ts, err := time.Parse(time.RFC3339, v.Attributes.UpdatedAt)
		if err != nil || ts.IsZero() {
			ts = time.Now()
		}

		fields := map[string]any{
			"vehicle_id": v.ID,
			"label":      v.Attributes.Label,
			"lat":        *v.Attributes.Latitude,
			"lng":        *v.Attributes.Longitude,
			"status":     v.Attributes.CurrentStatus,
		}

		if v.Attributes.Bearing != nil {
			fields["bearing"] = *v.Attributes.Bearing
		}

		if v.Attributes.Speed != nil {
			fields["speed"] = *v.Attributes.Speed
		}

		records = append(records, domain.RawRecord{
			SourceID:  v.ID,
			Timestamp: ts,
			Fields:    fields,
		})
	}

	return records, nil
}
