package handler

import (
	"github.com/labstack/echo/v4"
)

func (h *Handler) SetUserIDMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Request().Header.Get("X-Forwarded-User")
		c.Set("userID", userID)
		return next(c)
	}
}
