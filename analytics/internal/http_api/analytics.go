package http_handler

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *HttpAPI) getDevelopersReport(c echo.Context) error {
	profitLog, err := h.s.GetDevelopersReport(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError(err))
	}
	return c.JSON(http.StatusOK, ResponseOK(profitLog))
}

func (h *HttpAPI) getProfitReport(c echo.Context) error {
	profitLog, err := h.s.GetProfitReport(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError(err))
	}
	return c.JSON(http.StatusOK, ResponseOK(profitLog))
}

func (h *HttpAPI) getRevenueSourceReport(c echo.Context) error {
	profitLog, err := h.s.GetRevenueSourceReport(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError(err))
	}
	return c.JSON(http.StatusOK, ResponseOK(profitLog))
}
