package handler

import (
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"
)

type Character struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}

type GetCharactersResponse struct {
	Characters []Character `json:"characters"`
}

func (h *Handler) GetCharacters(c echo.Context) error {
	var characters []Character
	err := h.db.Select(&characters, "SELECT id, description, created_at FROM characters")
	if err != nil {
		slog.Error("Failed to fetch characters", slog.String("error", err.Error()))
		return c.JSON(500, map[string]string{"error": "Failed to fetch characters"})
	}

	res := GetCharactersResponse{
		Characters: characters,
	}
	return c.JSON(200, res)
}
