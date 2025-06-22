package handler

import (
	"database/sql"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Character struct {
	ID          string    `json:"id" db:"id"`
	Description string    `json:"description" db:"description"`
	IconURL     string    `json:"iconUrl" db:"icon_url"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
}

type GetCharactersResponse struct {
	Characters []Character `json:"characters"`
}

func (h *Handler) GetCharacters(c echo.Context) error {
	var characters []Character
	err := h.db.Select(&characters, "SELECT id, description, icon_url, created_at FROM characters")
	if err != nil {
		slog.Error("Failed to fetch characters", slog.String("error", err.Error()))
		return c.JSON(500, map[string]string{"error": "Failed to fetch characters"})
	}

	res := GetCharactersResponse{
		Characters: characters,
	}
	return c.JSON(200, res)
}

type CharacterIcon struct {
	Icon        []byte `db:"icon"`
	ContentType string `db:"content_type"`
}

func (h *Handler) GetCharacterIcon(c echo.Context) error {
	id := c.Param("id")

	var icon CharacterIcon
	err := h.db.Get(&icon, "SELECT icon, content_type FROM ai_icons WHERE character_id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(404, map[string]string{"error": "Character icon not found"})
		}
		slog.Error("Failed to fetch character icon", slog.String("error", err.Error()))
		return c.JSON(500, map[string]string{"error": "Failed to fetch character icon"})
	}

	c.Response().Header().Set("Content-Type", icon.ContentType)
	c.Response().Header().Set("Cache-Control", "public, max-age=31536000") // Cache for 1 year
	return c.Blob(200, icon.ContentType, icon.Icon)
}

type PostCharacterRequest struct {
	Instruction string `json:"instruction"`
	Description string `json:"description"`
}

func (h *Handler) PostCharacter(c echo.Context) error {
	var req PostCharacterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request format"})
	}

	id := uuid.New()
	iconURL := h.config.DefaultAIIconURL

	tx := h.db.MustBegin()
	defer tx.Rollback() // Ensure rollback on error

	_, err := tx.Exec("INSERT INTO characters (id, instruction, description, icon_url) VALUES (?, ?, ?, ?)", id, req.Instruction, req.Description, iconURL)
	if err != nil {
		slog.Error("Failed to save character", slog.String("error", err.Error()))
		return c.JSON(500, map[string]string{"error": "Failed to save character"})
	}

	if err := tx.Commit(); err != nil {
		slog.Error("Failed to commit transaction", slog.String("error", err.Error()))
		return c.JSON(500, map[string]string{"error": "Failed to commit transaction"})
	}

	response := map[string]string{
		"id": id.String(),
	}
	return c.JSON(200, response)
}

func (h *Handler) PutCharacterIcon(c echo.Context) error {
	id := c.Param("id")
	file, err := c.FormFile("icon")
	if err != nil {
		slog.Error("Failed to get form file", slog.String("error", err.Error()))
		return c.JSON(400, map[string]string{"error": "Invalid form file"})
	}
	image, err := file.Open()
	if err != nil {
		slog.Error("Failed to open form file", slog.String("error", err.Error()))
		return c.JSON(500, map[string]string{"error": "Failed to open form file"})
	}
	defer image.Close()
	blob, err := io.ReadAll(image)
	if err != nil {
		slog.Error("Failed to read form file", slog.String("error", err.Error()))
		return c.JSON(500, map[string]string{"error": "Failed to read form file"})
	}

	tx := h.db.MustBegin()
	defer tx.Rollback() // Ensure rollback on error

	contentType := http.DetectContentType(blob)
	_, err = tx.Exec("REPLACE ai_icons (character_id, icon, content_type) VALUES (?, ?, ?)",
		id,
		blob,
		contentType,
	)
	if err != nil {
		slog.Error("Failed to update character icon", slog.String("error", err.Error()))
		return c.JSON(500, map[string]string{"error": "Failed to update character icon"})
	}

	iconURL := "/api/characters/" + id + "/icon"
	_, err = tx.Exec("UPDATE characters SET icon_url = ? WHERE id = ?", iconURL, id)
	if err != nil {
		slog.Error("Failed to update character icon URL", slog.String("error", err.Error()))
		return c.JSON(500, map[string]string{"error": "Failed to update character icon URL"})
	}

	if err := tx.Commit(); err != nil {
		slog.Error("Failed to commit transaction", slog.String("error", err.Error()))
		return c.JSON(500, map[string]string{"error": "Failed to commit transaction"})
	}

	return c.NoContent(204)
}
