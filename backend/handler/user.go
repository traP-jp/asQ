package handler

import "github.com/labstack/echo/v4"

type User struct {
	ID string `json:"userId"`
}

func (h *Handler) GETMe(c echo.Context) error {
	userID := c.Request().Header.Get("X-Forwarded-User")
	return c.JSON(200, User{ID: userID})
}
