package http_handler

import (
	"async_course/auth"
	"log/slog"
	"net/http"
	"slices"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

const (
	RoleDeveloper  = "developer"
	RoleAdmin      = "admin"
	RoleManager    = "manager"
	RoleAccountant = "accountant"
)

func (h *HttpAPI) RegisterPublic(g *echo.Group) {
	g.GET("/status", h.hello)
	g.POST("/login", h.login)
}

func (h *HttpAPI) RegisterAPI(g *echo.Group) {
	g.POST("/create-account", h.requireRoles(h.createAccount, []string{RoleManager, RoleAdmin}))
	g.POST("/change-account-role", h.requireRoles(h.changeAccountRole, []string{RoleManager, RoleAdmin}))
}

func (h *HttpAPI) hello(c echo.Context) error {
	return c.JSON(http.StatusOK, ResponseOK(nil))
}

func (h *HttpAPI) requireRoles(fn echo.HandlerFunc, roles []string) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, ok := c.Get("user").(*jwt.Token)
		if !ok {
			slog.Error("could not find jwt token in request context", "error", auth.ErrTokenNotFound)
			return c.JSON(http.StatusInternalServerError, auth.ErrTokenNotFound)
		}

		claims, ok := token.Claims.(*auth.JwtCustomClaims)
		if !ok {
			slog.Error("cannot cast to *jwtClaims", "error", auth.ErrInvalidJwtClaimsFormat)
			return c.JSON(http.StatusForbidden, auth.ErrInvalidJwtClaimsFormat)
		}

		if claims == nil {
			slog.Error("empty claims in provided token", "error", auth.ErrInvalidJwtClaimsFormat)
			return c.JSON(http.StatusForbidden, auth.ErrInvalidJwtClaimsFormat)
		}

		role := claims.Role
		if !slices.Contains(roles, role) {
			return c.JSON(http.StatusForbidden, ResponseError(auth.ErrInsufficientPrivileges))
		}
		return fn(c)
	}

}
