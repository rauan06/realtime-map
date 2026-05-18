package controller

import (
	"context"
	"embed"
	"encoding/json"
	"io/fs"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/rauan06/realtime-map/dashboard/internal/hub"
	"github.com/rauan06/realtime-map/go-commons/pkg/logger"
)

const writeBufferSize = 1024

type wsMessage struct {
	Layer   string          `json:"layer"`
	Payload json.RawMessage `json:"payload"`
}

func Register(mux *http.ServeMux, h *hub.Hub, templates fs.FS, l logger.Interface) {
	mux.Handle("/", http.FileServer(http.FS(templates)))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/ws", wsHandler(h, l))
}

func wsHandler(h *hub.Hub, l logger.Interface) http.HandlerFunc {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  writeBufferSize,
		WriteBufferSize: writeBufferSize,
		CheckOrigin:     func(_ *http.Request) bool { return true },
	}

	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		ch := h.Subscribe()
		defer h.Unsubscribe(ch)

		// Reader goroutine: detect client disconnect & exit the writer loop.
		go func() {
			defer cancel()
			for {
				if _, _, err := conn.NextReader(); err != nil {
					return
				}
			}
		}()

		defer conn.Close()
		for {
			select {
			case <-ctx.Done():
				return
			case m, ok := <-ch:
				if !ok {
					return
				}
				envelope := wsMessage{Layer: m.Layer, Payload: m.Payload}
				data, err := json.Marshal(envelope)
				if err != nil {
					l.Error("dashboard ws encode: %v", err)
					continue
				}
				if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
					return
				}
			}
		}
	}
}

// EmbedFS is a convenience type so the app package can pass its //go:embed FS in.
type EmbedFS = embed.FS
