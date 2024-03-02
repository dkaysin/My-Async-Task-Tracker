package http_handler

import (
	"async_course/accounting"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *HttpAPI) RegisterPublic(g *echo.Group) {
	g.GET("/status", h.status)
}

func (h *HttpAPI) RegisterAPI(g *echo.Group) {
	g.GET("/balance-summary", h.requireRoles(h.getBalanceSummary, []string{accounting.RoleDeveloper}))
	g.GET("/balance-log", h.requireRoles(h.getBalanceLog, []string{accounting.RoleDeveloper}))
	g.GET("/profit-log", h.requireRoles(h.getProfitLog, []string{accounting.RoleManager, accounting.RoleAdmin, accounting.RoleAccountant}))
	g.POST("/close-balance", h.requireRoles(h.closeBalance, []string{accounting.RoleAdmin}))
}

func (h *HttpAPI) status(c echo.Context) error {
	return c.JSON(http.StatusOK, ResponseOK(nil))
}
