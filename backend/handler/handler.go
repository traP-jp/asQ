package handler

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/traP-jp/h25s_05/backend/llm"
)

type Handler struct {
	db     *sqlx.DB
	llmsvc *llm.Service
}

func NewHandler(db *sqlx.DB, llmsvc *llm.Service) *Handler {
	return &Handler{
		db:     db,
		llmsvc: llmsvc,
	}
}

func (h *Handler) SetUpRoutes(api *echo.Group) {
	api.GET("/ping", func(c echo.Context) error {
		return c.String(200, "pong")
	})

	api.GET("/sse/ai/:id", h.StreamAIResponse)

	api.GET("/chats", h.GETChats)
	api.POST("/chats", h.POSTChats)
	api.POST("/chats/:id/search", h.PostMessage)

	api.GET("/users/me", h.GETMe)
}
