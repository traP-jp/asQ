package handler

import (
	"log/slog"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traP-jp/h25s_05/backend/event"
	"github.com/traP-jp/h25s_05/backend/llm"
)

type Message struct {
	ID        string `json:"userId" db:"user_id"`
	Message   string `json:"message" db:"content"`
	CreatedAt string `json:"createdAt" db:"created_at"`
}

func (h *Handler) GETMessageID(c echo.Context) error {
	id := c.Param("id")
	var messages Message
	err := h.db.Get(&messages, "SELECT user_id, content, created_at FROM messages WHERE id = ?", id)
	if err != nil {
		return c.String(500, err.Error())
	}
	return c.JSON(200, messages)
}

type PostMessageRequest struct {
	Message     string `json:"message"`
	CharacterID string `json:"characterId"`
}

type PostMessageResponse struct {
	ID string `json:"id"`
}

func (h *Handler) PostMessage(c echo.Context) error {
	var req PostMessageRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request format"})
	}
	chatID := c.Param("id")
	messageID := uuid.New()
	userID := c.Get("userID").(string)

	_, err := h.db.Exec("INSERT INTO messages (id, chat_id, user_id, content) VALUES (?, ?, ?, ?)", messageID, chatID, userID, req.Message)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to save message"})
	}

	h.em.Publish(chatID, event.Event{Type: "message", Data: messageID})

	responseID, whenComplete := h.llmsvc.AskQuestion(req.Message, "", llm.MCP{
		ServerLabel: "deepwiki",
		ServerURL:   "https://mcp.deepwiki.com/mcp",
	}) // TODO: Implement character ID handling

	h.em.Publish(chatID, event.Event{Type: "response", Data: responseID})

	go func() {
		res := <-whenComplete
		if res.Err != nil {
			slog.Error("Failed to get response from LLM", slog.String("error", res.Err.Error()))
			return
		}
		_, err := h.db.Exec("INSERT INTO responses (id, openai_id, ai_id, chat_id, message_id, content) VALUES (?, ?, ?, ?, ?, ?)",
			responseID,
			res.ID,
			req.CharacterID,
			chatID, messageID,
			res.Text)
		if err != nil {
			slog.Error("Failed to save response", slog.String("error", err.Error()), slog.String("responseID", responseID.String()))
		}
	}()

	return c.JSON(200, PostMessageResponse{ID: responseID.String()})
}
