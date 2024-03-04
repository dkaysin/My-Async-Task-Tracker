package http_handler

import (
	"async_course/auth"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *HttpAPI) RegisterPublic(g *echo.Group) {
	g.GET("/status", h.hello)
	g.POST("/login", h.login)
}

func (h *HttpAPI) RegisterAPI(g *echo.Group) {
	g.POST("/create-account", h.requireRoles(h.createAccount, []string{auth.RoleManager, auth.RoleAdmin}))
	g.POST("/change-account-role", h.requireRoles(h.changeAccountRole, []string{auth.RoleManager, auth.RoleAdmin}))
}

func (h *HttpAPI) hello(c echo.Context) error {
	return c.JSON(http.StatusOK, ResponseOK(nil))
}
