package handler

import (
	"net/http"
	"strconv"

	"github.com/kenshin579/tutorials-go/rbac/backend/domain"
	"github.com/kenshin579/tutorials-go/rbac/backend/usecase"
	"github.com/labstack/echo/v4"
)

type OrderHandler struct {
	orderUsecase usecase.OrderUsecase
}

func NewOrderHandler(ou usecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{orderUsecase: ou}
}

type createOrderRequest struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

func (h *OrderHandler) Create(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	var req createOrderRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	order := &domain.Order{
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		OrderedBy: userID,
	}

	if err := h.orderUsecase.Create(order); err != nil {
		if err == usecase.ErrProductNotFound {
			return echo.NewHTTPError(http.StatusBadRequest, "product not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) List(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	roles := c.Get("roles").([]string)

	orders, err := h.orderUsecase.List(userID, roles)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) Get(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	order, err := h.orderUsecase.GetByID(uint(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "order not found")
	}

	return c.JSON(http.StatusOK, order)
}

type updateOrderStatusRequest struct {
	Status string `json:"status"`
}

func (h *OrderHandler) UpdateStatus(c echo.Context) error {
	roles := c.Get("roles").([]string)

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	var req updateOrderStatusRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	if err := h.orderUsecase.UpdateStatus(uint(id), domain.OrderStatus(req.Status), roles); err != nil {
		if err == usecase.ErrForbiddenTransition {
			return echo.NewHTTPError(http.StatusForbidden, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "status updated"})
}

func (h *OrderHandler) Cancel(c echo.Context) error {
	roles := c.Get("roles").([]string)

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	if err := h.orderUsecase.Cancel(uint(id), roles); err != nil {
		if err == usecase.ErrForbiddenTransition {
			return echo.NewHTTPError(http.StatusForbidden, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "order cancelled"})
}
