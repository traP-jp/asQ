package handler

import "github.com/labstack/echo/v4"

func (h *Handler) GETMe(c echo.Context) error {
	userID := c.Request().Header.Get("X-Forwarded-User")
	return c.JSON(200, userID)
}
