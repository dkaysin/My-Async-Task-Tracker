package http_handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *HttpAPI) RegisterAPI(g *echo.Group) {
	g.GET("/status", h.hello)
}

func (h *HttpAPI) hello(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}