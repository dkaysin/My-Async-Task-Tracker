package http_handler

import (
	"async_course/task"
	service "async_course/task/internal/service"
	"net/http"

	"github.com/go-playground/validator"
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
		return payload, task.ErrPayloadValidationFailed
	}
	validate := validator.New()
	if err := validate.Struct(payload); err != nil {
		return payload, task.ErrPayloadValidationFailed
	}
	return payload, nil
}

func ResponseOK(data interface{}) map[string]interface{} {
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
