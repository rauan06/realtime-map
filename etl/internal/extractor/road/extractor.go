package road

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/rauan06/realtime-map/etl/internal/domain"
)

const defaultAPIURL = "http://road-snowdrift-forecast:8080/api/v1/predictions/live"

type Extractor struct {
	client *http.Client
	url    string
}

func New(timeout time.Duration, endpoint string) *Extractor {
	if endpoint == "" {
		endpoint = defaultAPIURL
	}
	return &Extractor{
		client: &http.Client{Timeout: timeout},
		url:    endpoint,
	}
}

type apiResponse struct {
	GeneratedAt string         `json:"generated_at"`
	Predictions []apiPredition `json:"predictions"`
}

type apiPredition struct {
	ID              int64       `json:"id"`
	RoadID          string      `json:"road_id"`
	RoadName        string      `json:"road_name"`
	Segment         apiSegment  `json:"segment"`
	RestrictionType string      `json:"restriction_type"`
	Severity        string      `json:"severity"`
	Reason          []string    `json:"reason"`
	Confidence      float64     `json:"confidence"`
	Weather         apiWeather  `json:"weather_context"`
	PredictedStart  string      `json:"predicted_start"`
	PredictedEnd    string      `json:"predicted_end"`
	PredictedAt     string      `json:"predicted_at"`
}

type apiSegment struct {
	StartKM  int    `json:"start_km"`
	EndKM    int    `json:"end_km"`
	Geometry string `json:"geometry"`
}

type apiWeather struct {
	PrecipitationMM float64 `json:"precipitation_mm"`
	SnowfallMM      float64 `json:"snowfall_mm"`
	TemperatureC    float64 `json:"temperature_c"`
	TempMinC        float64 `json:"temp_min_c"`
	TempMaxC        float64 `json:"temp_max_c"`
	WindSpeedMS     float64 `json:"wind_speed_ms"`
	WindGustMS      float64 `json:"wind_gust_ms"`
	DewpointC       float64 `json:"dewpoint_c"`
}

func (e *Extractor) Extract(ctx context.Context) ([]domain.RawRecord, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, e.url, nil)
	if err != nil {
		return nil, fmt.Errorf("road extract: new request: %w", err)
	}

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("road extract: do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("road extract: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("road extract: read body: %w", err)
	}

	var data apiResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("road extract: unmarshal: %w", err)
	}

	records := make([]domain.RawRecord, 0, len(data.Predictions))
	for _, p := range data.Predictions {
		coords := parseMultiLineString(p.Segment.Geometry)
		if len(coords) == 0 {
			continue
		}

		ts := parseTime(p.PredictedAt)
		records = append(records, domain.RawRecord{
			SourceID:  p.RoadID,
			Timestamp: ts,
			Fields: map[string]interface{}{
				"prediction_id":    p.ID,
				"road_id":          p.RoadID,
				"road_name":        p.RoadName,
				"start_km":         p.Segment.StartKM,
				"end_km":           p.Segment.EndKM,
				"coords":           coords,
				"restriction_type": p.RestrictionType,
				"severity":         p.Severity,
				"reason":           p.Reason,
				"confidence":       p.Confidence,
				"weather":          p.Weather,
				"predicted_start":  p.PredictedStart,
				"predicted_end":    p.PredictedEnd,
			},
		})
	}

	return records, nil
}

// parseMultiLineString converts a Web-Mercator MULTILINESTRING wkt as emitted
// by road-snowdrift-forecast into a list of [lng, lat] WGS84 pairs (flattened,
// ignoring sub-line breaks — the dashboard renders them as a single polyline
// per road which is correct for the highway segments we visualise).
func parseMultiLineString(wkt string) [][2]float64 {
	const prefix = "MULTILINESTRING"
	wkt = strings.TrimSpace(wkt)
	if !strings.HasPrefix(strings.ToUpper(wkt), prefix) {
		return nil
	}
	wkt = wkt[len(prefix):]
	wkt = strings.TrimSpace(wkt)
	if strings.EqualFold(wkt, "EMPTY") {
		return nil
	}

	wkt = strings.TrimPrefix(wkt, "(")
	wkt = strings.TrimSuffix(wkt, ")")

	var out [][2]float64
	for _, group := range strings.Split(wkt, "),(") {
		group = strings.Trim(group, "() ")
		for _, pair := range strings.Split(group, ",") {
			parts := strings.Fields(strings.TrimSpace(pair))
			if len(parts) < 2 {
				continue
			}
			x, err1 := strconv.ParseFloat(parts[0], 64)
			y, err2 := strconv.ParseFloat(parts[1], 64)
			if err1 != nil || err2 != nil {
				continue
			}
			lng, lat := webMercatorToWGS84(x, y)
			out = append(out, [2]float64{lng, lat})
		}
	}
	return out
}

// webMercatorToWGS84 inverts the projection used in the live_predictions
// payload (EPSG:3857 → EPSG:4326).
func webMercatorToWGS84(x, y float64) (lng, lat float64) {
	const r = 6378137.0
	lng = (x / r) * 180.0 / math.Pi
	lat = (math.Atan(math.Exp(y/r))*2 - math.Pi/2) * 180.0 / math.Pi
	return lng, lat
}

func parseTime(s string) time.Time {
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t
	}
	return time.Now().UTC()
}
