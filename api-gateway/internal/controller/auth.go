package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/rauan06/realtime-map/api-gateway/internal/auth"
)

// LoginConfig is the dev-only login endpoint config. In production the
// shared secret would be replaced with proper device credential validation.
type LoginConfig struct {
	Issuer       *auth.Issue
	SharedSecret string // device must present this to log in
}

type loginRequest struct {
	DeviceID     string `json:"device_id"`
	SharedSecret string `json:"shared_secret"`
}

type loginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// RegisterAuth wires POST /auth/login. The handler issues an HS256 token to
// any device that presents the configured shared_secret — a deliberately
// minimal stand-in for an IdP, but enough to make /ws gated.
func RegisterAuth(mux *http.ServeMux, cfg LoginConfig) {
	mux.HandleFunc("POST /auth/login", func(w http.ResponseWriter, r *http.Request) {
		var req loginRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)

			return
		}

		if req.DeviceID == "" {
			http.Error(w, `{"error":"device_id required"}`, http.StatusBadRequest)

			return
		}

		if cfg.SharedSecret != "" && req.SharedSecret != cfg.SharedSecret {
			http.Error(w, `{"error":"invalid shared_secret"}`, http.StatusUnauthorized)

			return
		}

		token, exp, err := cfg.Issuer.Sign(req.DeviceID)
		if err != nil {
			http.Error(w, `{"error":"failed to sign token"}`, http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(loginResponse{Token: token, ExpiresAt: exp}); err != nil {
			// Already wrote headers; nothing else we can do.
			_ = err
		}
	})
}
