package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h *Handler) EventStream(c echo.Context) error {
	chatID := c.Param("id")

	w := c.Response()
	sse := StartSSE(w)

	err := sse.WriteMessage("start", "Event stream started")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to write SSE start message"})
	}

	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	stream := h.em.Subscribe(ctx, chatID)

	for event := range stream {
		var err error
		switch event.Type {
		case "message":
			err = sse.WriteMessage("user", event.Data.(uuid.UUID).String())
		case "response":
			err = sse.WriteMessage("ai", event.Data.(uuid.UUID).String())
		}
		if err != nil {
			slog.Error("Failed to write SSE message", "error", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to write SSE message"})
		}
	}

	if err := sse.WriteMessage("close", "Stream closed"); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to write SSE close message"})
	}

	return nil
}

type AIResponseMessage struct {
	Message string `json:"message"`
}

type AIResponseStartMessage struct {
	CharacterID uuid.UUID `json:"characterId"`
}

func (h *Handler) StreamAIResponse(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}

	w := c.Response()
	sse := StartSSE(w)

	err = sse.WriteJSON("start", AIResponseStartMessage{CharacterID: id})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to write SSE start message"})
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream, err := h.llmsvc.Subscribe(ctx, id)
	if err != nil {
		slog.Error("Failed to subscribe to AI response stream", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to subscribe to AI response stream"})
	}
	for data := range stream {
		err := sse.WriteJSON("data", AIResponseMessage{Message: data.TextDelta})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to write SSE data"})
		}
	}

	if err := sse.WriteMessage("close", "Stream closed"); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to write SSE close message"})
	}

	<-c.Request().Context().Done()

	return nil
}

type SSEWriter struct {
	w http.ResponseWriter
}

func StartSSE(w http.ResponseWriter) *SSEWriter {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	return &SSEWriter{w}
}

func (sse *SSEWriter) WriteJSON(event string, v any) error {
	if _, err := sse.w.Write([]byte("event: " + event + "\n")); err != nil {
		return err
	}
	if _, err := sse.w.Write([]byte("data: ")); err != nil {
		return err
	}
	if err := json.NewEncoder(sse.w).Encode(v); err != nil {
		return err
	}
	if _, err := sse.w.Write([]byte("\n\n")); err != nil {
		return err
	}
	sse.w.(http.Flusher).Flush()
	return nil
}

func (sse *SSEWriter) WriteMessage(event string, s string) error {
	if _, err := sse.w.Write([]byte("event: " + event + "\n")); err != nil {
		return err
	}
	if _, err := sse.w.Write([]byte("data: ")); err != nil {
		return err
	}
	if _, err := sse.w.Write([]byte(s)); err != nil {
		return err
	}
	sse.w.Write([]byte("\n\n"))
	sse.w.(http.Flusher).Flush()
	return nil
}
