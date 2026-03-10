package handler

import (
	"net/http"
	"strconv"

	"github.com/kenshin579/tutorials-go/rbac/backend/domain"
	"github.com/kenshin579/tutorials-go/rbac/backend/usecase"
	"github.com/labstack/echo/v4"
)

type ProductHandler struct {
	productUsecase usecase.ProductUsecase
}

func NewProductHandler(pu usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{productUsecase: pu}
}

func (h *ProductHandler) List(c echo.Context) error {
	roles := c.Get("roles").([]string)

	products, err := h.productUsecase.List(roles)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) Get(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	product, err := h.productUsecase.GetByID(uint(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "product not found")
	}

	return c.JSON(http.StatusOK, product)
}

type createProductRequest struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func (h *ProductHandler) Create(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	var req createProductRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	product := &domain.Product{
		Name:      req.Name,
		Price:     req.Price,
		Status:    domain.ProductStatusActive,
		CreatedBy: userID,
	}

	if err := h.productUsecase.Create(product); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, product)
}

type updateProductRequest struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func (h *ProductHandler) Update(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	product, err := h.productUsecase.GetByID(uint(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "product not found")
	}

	var req updateProductRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Price > 0 {
		product.Price = req.Price
	}

	if err := h.productUsecase.Update(product); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) Delete(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	if err := h.productUsecase.Delete(uint(id)); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "product deleted"})
}

type updateStatusRequest struct {
	Status string `json:"status"`
}

func (h *ProductHandler) UpdateStatus(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	product, err := h.productUsecase.GetByID(uint(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "product not found")
	}

	var req updateStatusRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	product.Status = domain.ProductStatus(req.Status)
	if err := h.productUsecase.Update(product); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, product)
}
