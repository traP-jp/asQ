package handler

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	db *sqlx.DB
}

func NewHandler(db *sqlx.DB) *Handler {
	return &Handler{
		db: db,
	}
}

func (h *Handler) SetUpRoutes(api *echo.Group) {
	api.Use(h.SetUserIDMiddleware)
	
	api.GET("/ping", func(c echo.Context) error {
		return c.String(200, "pong")
	})

	api.GET("/chats", h.GETChats)

	api.POST("/chats", h.POSTChats)

	api.GET("/users/me", h.GETMe)

	api.GET("chat/log", h.GETChatLog)
}
