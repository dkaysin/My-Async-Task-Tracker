package http_handler

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

type loginReq struct {
	UserID       string `json:"user_id"`
	PasswordHash string `json:"password_hash"`
}

type loginRes struct {
	Token string `json:"token"`
}

func (h *HttpAPI) login(c echo.Context) error {
	payload, err := validatePayload[loginReq](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError(err))
	}
	jwtToken, err := h.s.Login(context.Background(), payload.UserID, payload.PasswordHash)
	if err != nil {
		return c.JSON(http.StatusForbidden, ResponseError(err))
	}
	return c.JSON(http.StatusOK, ResponseOK(loginRes{jwtToken}))
}
