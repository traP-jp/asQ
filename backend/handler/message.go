package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/traP-jp/h25s_05/backend/llm"
)

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

	id := h.llmsvc.AskQuestion(req.Message, "", llm.MCP{
		ServerLabel: "deepwiki",
		ServerURL:   "https://mcp.deepwiki.com/mcp",
	}) // TODO: Implement character ID handling

	return c.JSON(200, PostMessageResponse{ID: id.String()})
}
