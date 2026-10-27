package hub

import (
	"sync"
)

type Hub struct {
	mu  	    sync.Mutex
	subscribers map[chan struct{}]bool
}

func New() *Hub {
	return &Hub{
		subscribers: make(map[chan struct{}]bool),
	}
}

// Subscribe registers a new client channel to receive broadcast events.
func (h *Hub) Subscribe() chan struct{} {
	ch := make(chan struct{}, 1)

	h.mu.Lock()
	defer h.mu.Unlock()

	h.subscribers[ch] = true
	return ch
}

// Unsubscribe removes a client channel from the subscriber set.
func (h *Hub) Unsubscribe(ch chan struct{}) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.subscribers, ch)
	close(ch)
}

// Publish sends a reload signal to each client in the subscriber set.
func (h *Hub) Publish() {
	h.mu.Lock()
	defer h.mu.Unlock()

	for ch := range h.subscribers {
		select {
		case ch <- struct{}{}:
			// signal this channel
		default:
			// client is not reading so move on
		}
	}
}

