package handler

import (
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/traP-jp/h25s_05/backend/event"
	"github.com/traP-jp/h25s_05/backend/llm"
)

type Config struct {
	DefaultAIIconURL string `json:"defaultAIIconUrl"`
}

type Handler struct {
	config Config
	db     *sqlx.DB
	llmsvc *llm.Service
	em     *event.Manager

	// chatBusy is a map to track busy chats to prevent multiple LLM requests for the same chat
	// chatID -> bool
	chatBusy sync.Map
}

func NewHandler(config Config, db *sqlx.DB, llmsvc *llm.Service) *Handler {
	return &Handler{
		config: config,
		db:     db,
		llmsvc: llmsvc,
		em:     &event.Manager{},
	}
}

func (h *Handler) SetUpRoutes(api *echo.Group) {
	api.Use(h.EnsureUserMiddleware)

	api.GET("/ping", func(c echo.Context) error {
		return c.String(200, "pong")
	})

	api.GET("/sse/events/:id", h.EventStream)
	api.GET("/sse/ai/:id", h.StreamAIResponse)

	api.GET("/chats", h.GETChats)
	api.POST("/chats", h.POSTChats)
	api.POST("/chats/:id/search", h.PostMessage)
	api.GET("/chats/:id/log", h.GETChatLog)

	api.GET("/users/me", h.GETMe)

	api.GET("/characters", h.GetCharacters)
	api.POST("/characters", h.PostCharacter, h.AdminMiddleware)
	api.GET("/characters/:id/icon", h.GetCharacterIcon)
	api.PUT("/characters/:id/icon", h.PutCharacterIcon, h.AdminMiddleware)

	api.GET("/messages/:id", h.GETMessageID)
}
