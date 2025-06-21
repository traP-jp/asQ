package handler

import "github.com/labstack/echo/v4"

type Handler struct{}

func (h *Handler) SetUpRoutes(api *echo.Group) {
	api.GET("/ping", func(c echo.Context) error {
		return c.String(200, "pong")
	})
}
