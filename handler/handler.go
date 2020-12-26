package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/amarjeetanandsingh/myRetail/product"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service product.ProductServer
}

func New(service product.ProductServer) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) InstallRoutes(e *echo.Echo) {
	e.GET("/products/{id}", h.handleGetProductByID)
	e.PUT("/products/{id}", h.handleUpdateProductPrice)
}

func (h *Handler) handleGetProductByID(ctx echo.Context) error {
	productIDStr := ctx.Param("id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "invalid product id")
	}

	productDetail, err := h.service.FindProductByID(productID)
	if errors.Is(err, product.ErrNotFound) {
		return ctx.JSON(http.StatusNotFound, err)
	}
	if err != nil {
		return ctx.JSON(http.StatusExpectationFailed, err)
	}
	return ctx.JSON(http.StatusOK, productDetail)
}

func (h *Handler) handleUpdateProductPrice(ctx echo.Context) error {
	productIDStr := ctx.Param("id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "invalid product id")
	}

	prodPrice := product.CurrentPrice{}
	if err := ctx.Bind(&prodPrice); err != nil {
		return ctx.JSON(http.StatusBadRequest, "invalid input")
	}

	if err := h.service.UpdateProductPrice(productID, prodPrice); err != nil {
		return ctx.JSON(http.StatusNotModified, err)
	}
	return ctx.JSON(http.StatusOK, prodPrice)
}
