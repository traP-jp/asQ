package handler

import (
	"database/sql"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traP-jp/h25s_05/backend/event"
)

type Message struct {
	UserID    string    `json:"userId" db:"user_id"`
	Message   string    `json:"message" db:"content"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
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

	var isLLMStarted bool
	if _, loaded := h.chatBusy.LoadOrStore(chatID, true); loaded {
		slog.Error("Chat is busy", slog.String("chatID", chatID))
		return c.JSON(400, map[string]string{"error": "Chat is busy, please try again later"})
	}
	defer func() {
		if !isLLMStarted {
			h.chatBusy.Delete(chatID) // Remove chat from busy state if LLM processing didn't start
		}
	}()

	tx, err := h.db.Beginx()
	if err != nil {
		slog.Error("Failed to begin transaction", slog.String("error", err.Error()))
		return c.JSON(500, map[string]string{"error": "Failed to begin transaction"})
	}
	defer tx.Rollback() // Ensure rollback on error

	_, err = tx.Exec("INSERT INTO messages (id, chat_id, user_id, content) VALUES (?, ?, ?, ?)", messageID, chatID, userID, req.Message)
	if err != nil {
		slog.Error("Failed to save message", slog.String("error", err.Error()), slog.String("messageID", messageID.String()))
		return c.JSON(500, map[string]string{"error": "Failed to save message"})
	}

	var instruction string
	err = tx.Get(&instruction, "SELECT instruction FROM characters WHERE id = ?", req.CharacterID)
	if err != nil {
		slog.Error("Failed to get instruction for character", slog.String("error", err.Error()))
		return c.JSON(500, map[string]string{"error": "Failed to retrieve character instruction"})
	}

	var previousID string
	err = tx.Get(&previousID, "SELECT external_id FROM responses WHERE chat_id = ? ORDER BY created_at DESC LIMIT 1", chatID)
	if errors.Is(err, sql.ErrNoRows) {
		previousID = ""
	} else if err != nil {
		slog.Error("Failed to get previous response ID", slog.String("error", err.Error()))
		return c.JSON(500, map[string]string{"error": "Failed to retrieve previous response ID"})
	}

	if err := tx.Commit(); err != nil {
		slog.Error("Failed to commit transaction", slog.String("error", err.Error()))
		return c.JSON(500, map[string]string{"error": "Failed to commit transaction"})
	}

	h.em.Publish(chatID, event.Event{Type: "message", Data: messageID})

	responseID, whenComplete := h.llmsvc.AskQuestion(req.Message, instruction, previousID)

	h.em.Publish(chatID, event.Event{Type: "response", Data: responseID})

	// Start the LLM processing
	slog.Info("LLM processing started", slog.String("chatID", chatID), slog.String("messageID", messageID.String()), slog.String("responseID", responseID.String()))
	isLLMStarted = true

	go func() {
		defer h.chatBusy.Delete(chatID) // Ensure chat is marked as not busy after processing
		res := <-whenComplete
		if res.Err != nil {
			slog.Error("Failed to get response from LLM", slog.String("error", res.Err.Error()))
			return
		}
		_, err := h.db.Exec("INSERT INTO responses (id, external_id, character_id, chat_id, message_id, content) VALUES (?, ?, ?, ?, ?, ?)",
			responseID,
			res.ExternalID,
			req.CharacterID,
			chatID,
			messageID,
			res.Text)
		if err != nil {
			slog.Error("Failed to save response", slog.String("error", err.Error()), slog.String("responseID", responseID.String()))
		}

		lastMessage := res.Text[max(0, strings.LastIndex(res.Text, "\n")):]
		_, err = h.db.Exec("UPDATE chats SET title = ? WHERE id = ?", lastMessage, chatID)
		if err != nil {
			slog.Error("Failed to update chat title", slog.String("error", err.Error()), slog.String("chatID", chatID))
			return
		}
	}()

	return c.JSON(200, PostMessageResponse{ID: responseID.String()})
}

type Response struct {
	Message     string `json:"message" db:"content"`
	CharacterID string `json:"characterId" db:"character_id"`
	CreatedAt   string `json:"createdAt" db:"created_at"`
}

type ChatLogResponse struct {
	Messages  []Message  `json:"messages"`
	Responses []Response `json:"responses"`
}

func (h *Handler) GETChatLog(c echo.Context) error {
	chatID := c.Param("id")
	var messages []Message
	err := h.db.Select(&messages, "SELECT user_id, content, created_at FROM messages WHERE chat_id = ?", chatID)
	if err != nil {
		return c.String(500, err.Error())
	}

	var responses []Response
	err = h.db.Select(&responses, "SELECT content, character_id, created_at FROM responses WHERE chat_id = ?", chatID)
	if err != nil {
		return c.String(500, err.Error())
	}

	chatLogResponse := ChatLogResponse{
		Messages:  messages,
		Responses: responses,
	}

	return c.JSON(200, chatLogResponse)
}
