package detector

import "math"

const degToRadFeat = math.Pi / 180.0

// FlightFeatures extracts a feature vector for an OpenSky-shaped flight
// payload: [velocity, altitude, sin(heading), cos(heading), on_ground].
// sin/cos handle the heading discontinuity at 360°.
func FlightFeatures(p map[string]any) ([]float64, float64, float64, string, bool) {
	id, ok := p["icao24"].(string)
	if !ok || id == "" {
		return nil, 0, 0, "", false
	}

	lat, lok := toFloat(p["lat"])
	lng, gok := toFloat(p["lng"])

	if !lok || !gok {
		return nil, 0, 0, "", false
	}

	velocity, _ := toFloat(p["velocity"])
	altitude, _ := toFloat(p["altitude"])
	heading, hok := toFloat(p["heading"])

	if !hok {
		heading = 0
	}

	onGround := boolFloat(p["on_ground"])

	vec := []float64{
		velocity,
		altitude,
		math.Sin(heading * degToRadFeat),
		math.Cos(heading * degToRadFeat),
		onGround,
	}

	return vec, lat, lng, id, true
}

// ShipFeatures extracts a feature vector for an AIS-shaped ship payload:
// [sog, sin(cog), cos(cog), sin(heading), cos(heading), nav_stat].
func ShipFeatures(p map[string]any) ([]float64, float64, float64, string, bool) {
	id, ok := p["mmsi"].(string)
	if !ok || id == "" {
		return nil, 0, 0, "", false
	}

	lat, lok := toFloat(p["lat"])
	lng, gok := toFloat(p["lng"])

	if !lok || !gok {
		return nil, 0, 0, "", false
	}

	sog, _ := toFloat(p["sog"])
	cog, _ := toFloat(p["cog"])
	heading, _ := toFloat(p["heading"])
	navStat, _ := toFloat(p["nav_stat"])

	vec := []float64{
		sog,
		math.Sin(cog * degToRadFeat),
		math.Cos(cog * degToRadFeat),
		math.Sin(heading * degToRadFeat),
		math.Cos(heading * degToRadFeat),
		navStat,
	}

	return vec, lat, lng, id, true
}

func toFloat(v any) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	default:
		return 0, false
	}
}

func boolFloat(v any) float64 {
	if b, ok := v.(bool); ok && b {
		return 1
	}

	return 0
}
