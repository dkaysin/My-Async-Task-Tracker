package http_handler

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

// create account

type createAccountReq struct {
	Name         string `json:"name"`
	PasswordHash string `json:"password_hash"`
	Role         string `json:"role"`
}

type createAccountRes struct {
	UserID string `json:"user_id"`
}

func (h *HttpAPI) createAccount(c echo.Context) error {
	payload, err := validatePayload[createAccountReq](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError(err))
	}
	userID, err := h.s.CreateAccount(context.Background(), payload.Name, payload.PasswordHash, payload.Role)
	if err != nil {
		return c.JSON(http.StatusForbidden, ResponseError(err))
	}
	return c.JSON(http.StatusOK, ResponseOK(createAccountRes{userID}))
}

// change account role

type changeAccountRoleReq struct {
	UserID  string `json:"user_id"`
	NewRole string `json:"new_role"`
}

func (h *HttpAPI) changeAccountRole(c echo.Context) error {
	payload, err := validatePayload[changeAccountRoleReq](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError(err))
	}
	err = h.s.ChangeAccountRole(context.Background(), payload.UserID, payload.NewRole)
	if err != nil {
		return c.JSON(http.StatusForbidden, ResponseError(err))
	}
	return c.JSON(http.StatusOK, ResponseOK(nil))
}
