package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/types/known/timestamppb"

	routepb "github.com/rauan06/realtime-map/go-commons/gen/proto/route"
)

const (
	bufferSize = 1024
)

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

func wsHandler(w http.ResponseWriter, r *http.Request, client routepb.RouteClient) {
	conn, err := websocketConn(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	ctx := r.Context()

	session, err := client.StartSession(ctx, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	err = conn.WriteMessage(websocket.TextMessage, []byte("{\"session_id\":\""+session.SessionId+"\"}"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	stream, err := client.RouteChat(ctx)
	if err != nil {
		err = conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"stream open failed"}`))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		return
	}

	defer func() {
		if _, err := stream.CloseAndRecv(); err != nil {
			log.Println(err)
		}
	}()

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			return
		}

		var in wsOBU
		if err := json.Unmarshal(data, &in); err != nil {
			err = conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"invalid payload"}`))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)

				return
			}

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
			err = conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"send failed"}`))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			return
		}
	}
}

func websocketConn(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  bufferSize,
		WriteBufferSize: bufferSize,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return conn, nil
}
