package event

import (
	"context"
	"sync"
)

type Event struct {
	Type string
	Data any
}

type subscriber struct {
	key string
	ch  chan Event
}

type Manager struct {
	mu          sync.Mutex
	subscribers []subscriber
}

func (m *Manager) Publish(key string, event Event) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, subscriber := range m.subscribers {
		if subscriber.key != key {
			continue // Skip subscribers that are not interested in this key
		}
		select {
		case subscriber.ch <- event:
			// Event sent successfully
		default:
			// Subscriber channel is full, skip sending to avoid blocking
		}
	}
}

func (m *Manager) Subscribe(ctx context.Context, key string) chan Event {
	m.mu.Lock()
	defer m.mu.Unlock()

	subscriberCh := make(chan Event, 100) // Buffered channel to avoid blocking
	m.subscribers = append(m.subscribers, subscriber{key, subscriberCh})

	go func() {
		<-ctx.Done() // Wait for context cancellation
		m.mu.Lock()
		defer m.mu.Unlock()

		// Remove the subscriber from the list
		for i, sub := range m.subscribers {
			if sub.ch == subscriberCh {
				m.subscribers = append(m.subscribers[:i], m.subscribers[i+1:]...)
				close(subscriberCh) // Close the channel to signal completion
				break
			}
		}
	}()

	return subscriberCh
}
