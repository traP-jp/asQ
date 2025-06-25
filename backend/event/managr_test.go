package event_test

import (
	"testing"

	"github.com/traP-jp/h25s_05/backend/event"
)

func TestPubSub_SameKey(t *testing.T) {
	t.Parallel()

	em := event.Manager{}

	subscriber := em.Subscribe(t.Context(), "test-key")

	em.Publish("test-key", event.Event{
		Type: "test-event",
	})

	select {
	case e := <-subscriber:
		if e.Type != "test-event" {
			t.Errorf("expected event type 'test-event', got '%s'", e.Type)
		}
	default:
		t.Error("expected to receive an event, but got none")
	}
}

func TestPubSub_DifferentKeys(t *testing.T) {
	t.Parallel()

	em := event.Manager{}

	subscriber := em.Subscribe(t.Context(), "test-key")

	em.Publish("different-key", event.Event{
		Type: "test-event",
	})

	select {
	case e := <-subscriber:
		t.Errorf("expected no event, but got '%s'", e.Type)
	default:
		// No event received, which is expected
	}
}
