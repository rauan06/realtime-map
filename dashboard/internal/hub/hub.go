package hub

import (
	"sync"
)

// Message is a single fan-out payload. Layer is the logical map layer
// ("flight"/"ship"/"transport"/"road"/"obu"); Payload is the raw JSON the
// browser will dispatch to its renderer for that layer.
type Message struct {
	Layer   string `json:"layer"`
	Payload []byte `json:"payload"`
}

// Hub is a tiny in-memory pub/sub that broadcasts every incoming Message to
// every connected subscriber. It is deliberately not back-pressure aware:
// slow subscribers drop the message rather than block the publisher.
type Hub struct {
	mu     sync.RWMutex
	subs   map[chan Message]struct{}
	bufLen int
}

const defaultSubBuffer = 64

func New(subBufferLen int) *Hub {
	if subBufferLen <= 0 {
		subBufferLen = defaultSubBuffer
	}

	return &Hub{
		subs:   make(map[chan Message]struct{}),
		bufLen: subBufferLen,
	}
}

func (h *Hub) Subscribe() chan Message {
	ch := make(chan Message, h.bufLen)

	h.mu.Lock()
	h.subs[ch] = struct{}{}
	h.mu.Unlock()

	return ch
}

func (h *Hub) Unsubscribe(ch chan Message) {
	h.mu.Lock()
	if _, ok := h.subs[ch]; ok {
		delete(h.subs, ch)
		close(ch)
	}
	h.mu.Unlock()
}

func (h *Hub) Publish(m Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for ch := range h.subs {
		select {
		case ch <- m:
		default:
			// drop on slow consumer — the realtime map prefers fresh updates
			// over a backlog of stale positions.
		}
	}
}

func (h *Hub) SubscriberCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return len(h.subs)
}
