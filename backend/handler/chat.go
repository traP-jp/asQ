package handler

import (
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Chat struct {
	ID                  string    `json:"id" db:"id"`
	ParticipantsUserIDs []string  `json:"userIds" db:"-"`
	Title               string    `json:"title" db:"title"`
	CreatedAt           time.Time `json:"createdAt" db:"created_at"`
}

type ChatsResponse struct {
	Chats []Chat `json:"chats"`
}

func (h *Handler) GETChats(c echo.Context) error {
	var chats []Chat
	err := h.db.Select(&chats, "SELECT id, title, created_at FROM chats")
	if err != nil {
		slog.Error("Failed to fetch chats", slog.String("error", err.Error()))
		return c.JSON(500, map[string]string{"error": "Failed to fetch chats"})
	}

	type ChatParticipant struct {
		ChatID string `db:"chat_id"`
		UserID string `db:"user_id"`
	}
	var chatParticipants []ChatParticipant
	err = h.db.Select(&chatParticipants, "SELECT DISTINCT chat_id, user_id FROM messages")
	if err != nil {
		slog.Error("Failed to fetch chat participants", slog.String("error", err.Error()))
		return c.JSON(500, map[string]string{"error": "Failed to fetch chat participants"})
	}

	chatsMap := make(map[string]*Chat)
	for i := range chats {
		chatsMap[chats[i].ID] = &chats[i]
	}

	for _, p := range chatParticipants {
		if chat, exists := chatsMap[p.ChatID]; exists {
			chat.ParticipantsUserIDs = append(chat.ParticipantsUserIDs, p.UserID)
		}
	}

	res := ChatsResponse{
		Chats: chats,
	}
	return c.JSON(200, res)
}

type PostChatsResponse struct {
	ID string `json:"id"`
}

func (h *Handler) POSTChats(c echo.Context) error {
	id := uuid.NewString()
	creatorID := c.Get("userID").(string)
	_, err := h.db.Exec("INSERT INTO chats (id, creator_id, title) VALUES (?, ?, ?)", id, creatorID, "New Chat")
	if err != nil {
		c.String(500, err.Error())
	}
	return c.JSON(200, PostChatsResponse{ID: id})
}

type ChatLog struct {
	ID        string `json:"id" db:"id"`
	ChatID    string `json:"chatId" db:"chat_id"`
	CreatedAt string `json:"createdAt" db:"created_at"`
}
