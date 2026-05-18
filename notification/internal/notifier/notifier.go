package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/rauan06/realtime-map/go-commons/pkg/logger"
)

const (
	webhookTimeout    = 5 * time.Second
	httpClientErrCode = 300
)

// Alert is the canonical event we emit. Event is one of "enter" / "exit".
type Alert struct {
	Event    string    `json:"event"`
	Layer    string    `json:"layer"`
	SourceID string    `json:"source_id"`
	Fence    string    `json:"fence"`
	Lat      float64   `json:"lat"`
	Lng      float64   `json:"lng"`
	At       time.Time `json:"at"`
}

// Notifier dispatches alerts. The Webhook field is optional; when empty,
// alerts are only logged. The HTTP client times out aggressively so a slow
// webhook never blocks the kafka consume loop.
type Notifier struct {
	Webhook string
	client  *http.Client
	l       logger.Interface
}

func New(webhook string, l logger.Interface) *Notifier {
	return &Notifier{
		Webhook: webhook,
		client:  &http.Client{Timeout: webhookTimeout},
		l:       l,
	}
}

func (n *Notifier) Dispatch(ctx context.Context, a Alert) {
	n.l.Info("ALERT %s layer=%s source=%s fence=%s at=(%.5f,%.5f)",
		a.Event, a.Layer, a.SourceID, a.Fence, a.Lat, a.Lng)

	if n.Webhook == "" {
		return
	}

	body, err := json.Marshal(a)
	if err != nil {
		n.l.Error("notifier: marshal alert: %v", err)

		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, n.Webhook, bytes.NewReader(body))
	if err != nil {
		n.l.Error("notifier: new request: %v", err)

		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := n.client.Do(req)
	if err != nil {
		n.l.Error("notifier: webhook post: %v", err)

		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= httpClientErrCode {
		n.l.Error("notifier: webhook returned status %d", resp.StatusCode)

		return
	}
}
