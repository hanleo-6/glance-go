package hub

import (
	"testing"
	"time"
)

func TestHub_Broadcast(t *testing.T) {
	h := New()

	const numSubscribers = 3
	subs := make([]chan struct{}, numSubscribers)
	for i := 0; i < numSubscribers; i++ {
		subs[i] = h.Subscribe()
	}

	// 1. Test that all N subscribers receive the event
	t.Run("all subscribers receive event", func(t *testing.T) {
		h.Publish()

		for i, ch := range subs {
			select {
			case <-ch:
				// success
			case <-time.After(50 * time.Millisecond):
				t.Fatalf("subscriber %d timed out waiting for event", i)
			}
		}
	})

	// 2. Test unsubscribe behavior
	t.Run("unsubscribed client does not receive events", func(t *testing.T) {
		unsubscribedCh := subs[0]
		h.Unsubscribe(unsubscribedCh)

		// Check unsubscribed channel is closed
		_, ok := <-unsubscribedCh
		if ok {
			t.Errorf("expected unsubscribed channel to be closed, but it was open")
		}

		// Check unsubscribed channel is removed from subscribers
		h.mu.Lock()
		_, exists := h.subscribers[unsubscribedCh]
		h.mu.Unlock()

		if exists {
			t.Errorf("unsubscribed channel still present in subscribers set")
		}

		// Check map length is now N-1
		h.mu.Lock()
		count := len(h.subscribers)
		h.mu.Unlock()

		if count != numSubscribers-1 {
			t.Errorf("expected %d subscribers after unsubscribe, got %d",
				numSubscribers-1, count)
		}

		// Test remaining subscribers receive events
		h.Publish()

		for i := 1; i < numSubscribers; i++ {
			select {
			case <-subs[i]:
				// success for this subscriber
			case <-time.After(50 * time.Millisecond):
				t.Fatalf("subscriber %d timed out waiting for event", i)
			}
		}
	})
}
