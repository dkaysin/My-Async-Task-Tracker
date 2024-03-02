package http_handler

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *HttpAPI) getBalanceSummary(c echo.Context) error {
	claims, err := getClaimsFromContext(c)
	if err != nil {
		return err
	}
	transactionsSummary, err := h.s.GetBalanceSummary(context.Background(), claims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError(err))
	}
	return c.JSON(http.StatusOK, ResponseOK(transactionsSummary))
}

func (h *HttpAPI) getBalanceLog(c echo.Context) error {
	claims, err := getClaimsFromContext(c)
	if err != nil {
		return err
	}
	transactionsLog, err := h.s.GetBalanceLog(context.Background(), claims.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError(err))
	}
	return c.JSON(http.StatusOK, ResponseOK(transactionsLog))
}

func (h *HttpAPI) getProfitLog(c echo.Context) error {
	profitLog, err := h.s.GetProfitLog(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError(err))
	}
	return c.JSON(http.StatusOK, ResponseOK(profitLog))
}

func (h *HttpAPI) closeBalance(c echo.Context) error {
	err := h.s.InitiateBalanceClose(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError(err))
	}
	return c.JSON(http.StatusOK, ResponseOK(nil))
}
