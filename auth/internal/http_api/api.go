package http_handler

import (
	"async_course/auth"
	service "async_course/auth/internal/service"
	"log/slog"
	"net/http"
	"slices"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

type HttpAPI struct {
	config *viper.Viper
	s      *service.Service
}

func NewHttpAPI(config *viper.Viper, s *service.Service) *HttpAPI {
	return &HttpAPI{
		config: config,
		s:      s,
	}
}

func validatePayload[T any](c echo.Context) (T, error) {
	var payload T
	if err := c.Bind(&payload); err != nil {
		return payload, auth.ErrPayloadValidationFailed
	}
	validate := validator.New()
	if err := validate.Struct(payload); err != nil {
		return payload, auth.ErrPayloadValidationFailed
	}
	return payload, nil
}

func ResponseOK(v interface{}) map[string]interface{} {
	return map[string]interface{}{
		"status": "ok",
		"data":   data,
	}
}

func ResponseError(err error) map[string]interface{} {
	return map[string]interface{}{
		"status": "error",
		"error":  err.Error(),
	}
}

func JwtMiddlewareErrorHandler(c echo.Context, err error) error {
	return c.JSON(http.StatusForbidden, ResponseError(err))
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
