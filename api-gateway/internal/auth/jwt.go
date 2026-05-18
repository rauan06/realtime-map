// Package auth implements HS256-signed JWTs for device authentication. Tokens
// carry the device id as the JWT subject and the issuer "map-api-gateway".
//
// This is intentionally minimal: no refresh tokens, no key rotation, single
// shared secret loaded from config. Sufficient for the realtime map's device
// auth (roadmap task 7) — replace with an external IdP for production.
package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	Issuer        = "map-api-gateway"
	DefaultTTL    = 24 * time.Hour
	bearerPrefix  = "Bearer "
	contextKeyDID ctxKey = "device_id"
)

type ctxKey string

// Issuer signs new tokens for a device.
type Issue struct {
	Secret []byte
	TTL    time.Duration
}

func NewIssue(secret string, ttl time.Duration) *Issue {
	if ttl <= 0 {
		ttl = DefaultTTL
	}
	return &Issue{Secret: []byte(secret), TTL: ttl}
}

func (i *Issue) Sign(deviceID string) (string, time.Time, error) {
	exp := time.Now().Add(i.TTL)
	claims := jwt.RegisteredClaims{
		Issuer:    Issuer,
		Subject:   deviceID,
		ExpiresAt: jwt.NewNumericDate(exp),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := tok.SignedString(i.Secret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("sign: %w", err)
	}
	return s, exp, nil
}

// Verify parses + validates a token string, returning the device id (subject).
func Verify(secret []byte, token string) (string, error) {
	parsed, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := parsed.Claims.(*jwt.RegisteredClaims)
	if !ok || !parsed.Valid {
		return "", errors.New("invalid token")
	}
	if claims.Issuer != Issuer {
		return "", errors.New("issuer mismatch")
	}
	if claims.Subject == "" {
		return "", errors.New("missing subject (device_id)")
	}
	return claims.Subject, nil
}

// Middleware gates the wrapped handler with bearer-token auth. The verified
// device id is stashed in the request context (use DeviceIDFromContext to
// read it back from downstream handlers).
//
// Skip auth on requests for which the token is unverifiable: when secret is
// nil, the middleware is a no-op so local dev can run with AUTH_ENABLED=false.
func Middleware(secret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if len(secret) == 0 {
			return next
		}
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tok := bearerFromRequest(r)
			if tok == "" {
				http.Error(w, `{"error":"missing bearer token"}`, http.StatusUnauthorized)
				return
			}
			deviceID, err := Verify(secret, tok)
			if err != nil {
				http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), contextKeyDID, deviceID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func bearerFromRequest(r *http.Request) string {
	if h := r.Header.Get("Authorization"); strings.HasPrefix(h, bearerPrefix) {
		return strings.TrimPrefix(h, bearerPrefix)
	}
	// Fallback: ?token= query param so browser WebSocket clients can authenticate
	// without sending custom headers.
	if v := r.URL.Query().Get("token"); v != "" {
		return v
	}
	return ""
}

func DeviceIDFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(contextKeyDID).(string)
	return v, ok
}
