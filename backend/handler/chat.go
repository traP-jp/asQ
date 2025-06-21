package handler

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Chat struct {
	ID    string `json:"id" db:"id"`
	Title string `json:"title" db:"title"`
}

type ChatsResponse struct {
	Chats []Chat `json:"chats"`
}

func (h *Handler) GETChats(c echo.Context) error {
	var chats []Chat
	err := h.db.Select(&chats, "SELECT id ,title FROM chats")
	if err != nil {
		c.String(500, err.Error())
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
	_, err := h.db.Exec("INSERT INTO chats (id, creator_id, title) VALUES (?, ?, ?)", id, creatorID, "aaaaa")
	if err != nil {
		c.String(500, err.Error())
	}
	return c.JSON(200, PostChatsResponse{ID: id})
}

func (h *Handler) GETChatLog(c echo.Context) error {

}
