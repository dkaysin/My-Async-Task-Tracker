package http_handler

import (
	"async_course/analytics"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *HttpAPI) RegisterPublic(g *echo.Group) {
	g.GET("/status", h.status)
}

func (h *HttpAPI) RegisterAPI(g *echo.Group) {
	g.GET("/developers-report", h.requireRoles(h.getDevelopersReport, []string{analytics.RoleAdmin}))
	g.GET("/profit-report", h.requireRoles(h.getProfitReport, []string{analytics.RoleAdmin}))
	g.GET("/revenue-source-report", h.requireRoles(h.getRevenueSourceReport, []string{analytics.RoleAdmin}))
}

func (h *HttpAPI) status(c echo.Context) error {
	return c.JSON(http.StatusOK, ResponseOK(nil))
}
