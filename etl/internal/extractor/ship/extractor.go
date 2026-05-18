package ship

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/rauan06/realtime-map/etl/internal/domain"
)

const apiURL = "https://meri.digitraffic.fi/api/ais/v1/locations"

type Extractor struct {
	client *http.Client
}

func New(timeout time.Duration) *Extractor {
	return &Extractor{
		client: &http.Client{Timeout: timeout},
	}
}

type featureCollection struct {
	Features []feature `json:"features"`
}

type feature struct {
	Geometry   geometry   `json:"geometry"`
	Properties properties `json:"properties"`
}

type geometry struct {
	Coordinates []float64 `json:"coordinates"`
}

type properties struct {
	MMSI      int     `json:"mmsi"`
	SOG       float64 `json:"sog"`
	COG       float64 `json:"cog"`
	Heading   float64 `json:"heading"`
	NavStat   int     `json:"navStat"`
	Timestamp int64   `json:"timestampExternal"`
}

func (e *Extractor) Extract(ctx context.Context) ([]domain.RawRecord, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("ship extract: new request: %w", err)
	}

	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Digitraffic-User", "realtime-map/1.0")

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ship extract: do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ship extract: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ship extract: read body: %w", err)
	}

	var data featureCollection
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("ship extract: unmarshal: %w", err)
	}

	records := make([]domain.RawRecord, 0, len(data.Features))
	for _, f := range data.Features {
		if len(f.Geometry.Coordinates) < 2 {
			continue
		}

		mmsi := strconv.Itoa(f.Properties.MMSI)
		lon := f.Geometry.Coordinates[0]
		lat := f.Geometry.Coordinates[1]

		ts := time.Unix(f.Properties.Timestamp/1000, 0)
		if f.Properties.Timestamp == 0 {
			ts = time.Now()
		}

		records = append(records, domain.RawRecord{
			SourceID:  mmsi,
			Timestamp: ts,
			Fields: map[string]interface{}{
				"mmsi":    mmsi,
				"lat":     lat,
				"lng":     lon,
				"sog":     f.Properties.SOG,
				"cog":     f.Properties.COG,
				"heading": f.Properties.Heading,
				"nav_stat": f.Properties.NavStat,
			},
		})
	}

	return records, nil
}
