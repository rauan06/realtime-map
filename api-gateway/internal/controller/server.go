package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	routepb "github.com/rauan06/realtime-map/go-commons/gen/proto/route"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type wsOBU struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp int64   `json:"timestamp,omitempty"`
}

// Register wires HTTP routes into the provided mux.
func Register(mux *http.ServeMux, routeClient routepb.RouteClient) {
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsHandler(w, r, routeClient)
	})
}

func wsHandler(w http.ResponseWriter, r *http.Request, routeClient routepb.RouteClient) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "websocket upgrade failed", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	ctx := r.Context()

	// Start session with producer
	session, err := routeClient.StartSession(ctx, nil)
	if err != nil {
		_ = conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"start session failed"}`))
		return
	}

	// Send session id to client
	_ = conn.WriteMessage(websocket.TextMessage, []byte("{\"session_id\":\""+session.SessionId+"\"}"))

	// Open client-streaming RPC
	stream, err := routeClient.RouteChat(ctx)
	if err != nil {
		_ = conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"stream open failed"}`))
		return
	}
	defer func() { _, _ = stream.CloseAndRecv() }()

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			return
		}

		var in wsOBU
		if err := json.Unmarshal(data, &in); err != nil {
			_ = conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"invalid payload"}`))
			continue
		}

		ts := time.Now()
		if in.Timestamp > 0 {
			ts = time.Unix(in.Timestamp, 0)
		}

		msg := &routepb.OBUData{
			SessionId: session.SessionId,
			Latitude:  in.Latitude,
			Longitude: in.Longitude,
			Timestamp: timestamppb.New(ts),
		}

		if err := stream.Send(msg); err != nil {
			_ = conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"send failed"}`))
			return
		}
	}
}
