package http_handler

import (
	"async_course/auth"
	service "async_course/auth/internal/service"

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
		"data":   v,
	}
}

func ResponseError(err error) map[string]interface{} {
	return map[string]interface{}{
		"status": "error",
		"error":  err,
	}
}
