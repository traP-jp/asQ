package handler

import (
	"github.com/labstack/echo/v4"
)

func (h *Handler) EnsureUserMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Request().Header.Get("X-Forwarded-User")
		_, err := h.db.Exec("INSERT INTO chats (id, username) VALUES (?, ?) ON DUPLICATE KEY UPDATE id = id", userID, userID)
		if err != nil {
			c.String(500, err.Error())
		}
		c.Set("userID", userID)
		return next(c)
	}
}
