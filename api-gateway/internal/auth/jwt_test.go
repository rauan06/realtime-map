package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestIssueAndVerify_RoundTrip(t *testing.T) {
	iss := NewIssue("test-secret", time.Hour)
	tok, exp, err := iss.Sign("device-123")
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	if tok == "" {
		t.Fatal("empty token")
	}
	if time.Until(exp) <= 0 {
		t.Fatalf("expiry not in future: %v", exp)
	}

	deviceID, err := Verify([]byte("test-secret"), tok)
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if deviceID != "device-123" {
		t.Errorf("device id: got %q want %q", deviceID, "device-123")
	}
}

func TestVerify_RejectsWrongSecret(t *testing.T) {
	iss := NewIssue("right-secret", time.Hour)
	tok, _, err := iss.Sign("device-abc")
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	if _, err := Verify([]byte("wrong-secret"), tok); err == nil {
		t.Error("expected signature error, got nil")
	}
}

func TestVerify_RejectsExpiredToken(t *testing.T) {
	// Sign a token with a past ExpiresAt directly — Issuer.Sign defaults
	// non-positive TTLs to 24h so we craft the claims manually here.
	past := time.Now().Add(-time.Hour)
	claims := jwt.RegisteredClaims{
		Issuer:    Issuer,
		Subject:   "device-x",
		ExpiresAt: jwt.NewNumericDate(past),
		IssuedAt:  jwt.NewNumericDate(past.Add(-time.Hour)),
	}
	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("test-secret"))
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	if _, err := Verify([]byte("test-secret"), signed); err == nil {
		t.Error("expected expired error, got nil")
	}
}

func TestMiddleware_AllowsValidBearer(t *testing.T) {
	iss := NewIssue("secret", time.Hour)
	tok, _, _ := iss.Sign("device-1")

	called := false
	h := Middleware([]byte("secret"))(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		did, ok := DeviceIDFromContext(r.Context())
		if !ok || did != "device-1" {
			t.Errorf("device id in context: got %q ok=%v", did, ok)
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/ws", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if !called {
		t.Error("handler not called for valid token")
	}
	if rec.Code != http.StatusOK {
		t.Errorf("status: got %d want 200", rec.Code)
	}
}

func TestMiddleware_RejectsMissingAndBadTokens(t *testing.T) {
	h := Middleware([]byte("secret"))(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	}))

	// Missing.
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/ws", nil))
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("missing token: status %d want 401", rec.Code)
	}

	// Bad.
	req := httptest.NewRequest(http.MethodGet, "/ws", nil)
	req.Header.Set("Authorization", "Bearer not-a-real-jwt")
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("bad token: status %d want 401", rec.Code)
	}
}

func TestMiddleware_DisabledIfNoSecret(t *testing.T) {
	h := Middleware(nil)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/ws", nil))
	if rec.Code != http.StatusTeapot {
		t.Errorf("no-secret middleware should pass-through; got %d", rec.Code)
	}
}

func TestMiddleware_AcceptsTokenViaQueryParam(t *testing.T) {
	iss := NewIssue("secret", time.Hour)
	tok, _, _ := iss.Sign("device-q")

	h := Middleware([]byte("secret"))(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/ws?token="+tok, nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("query param token rejected: status %d", rec.Code)
	}
}
