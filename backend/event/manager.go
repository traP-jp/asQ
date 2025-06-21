package event

import (
	"context"
	"sync"
)

type Event struct {
	Type string
	Data any
}

type Manager struct {
	mu          sync.Mutex
	subscribers []chan Event
}

func (m *Manager) Publish(key string, event Event) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, subscriber := range m.subscribers {
		select {
		case subscriber <- event:
			// Event sent successfully
		default:
			// Subscriber channel is full, skip sending to avoid blocking
		}
	}
}

func (m *Manager) Subscribe(ctx context.Context, key string) chan Event {
	m.mu.Lock()
	defer m.mu.Unlock()

	subscriber := make(chan Event, 100) // Buffered channel to avoid blocking
	m.subscribers = append(m.subscribers, subscriber)

	go func() {
		<-ctx.Done() // Wait for context cancellation
		m.mu.Lock()
		defer m.mu.Unlock()

		// Remove the subscriber from the list
		for i, sub := range m.subscribers {
			if sub == subscriber {
				m.subscribers = append(m.subscribers[:i], m.subscribers[i+1:]...)
				close(subscriber) // Close the channel to signal completion
				break
			}
		}
	}()

	return subscriber
}
