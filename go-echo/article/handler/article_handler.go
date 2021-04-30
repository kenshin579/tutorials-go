package handler

import (
	"net/http"

	"github.com/kenshin579/tutorials-go/go-echo/article/model"
	"github.com/kenshin579/tutorials-go/go-echo/article/utils"
	"github.com/labstack/echo"
)

type Handler struct {
	articleUsecase model.ArticleUsecase
}

func NewHandler(au model.ArticleUsecase) *Handler {
	return &Handler{
		articleUsecase: au,
	}
}

func (h *Handler) CreateArticle(c echo.Context) error {
	request := &model.ArticleRequest{}

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewValidatorError(err))
	}

	err := h.articleUsecase.CreateArticle(request)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, request)
}

func (h *Handler) GetArticle(c echo.Context) error {
	response, err := h.articleUsecase.GetArticle(c.Param("articleId"))

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, utils.NewError(err))
	}
	return c.JSON(http.StatusOK, response)
}

func (h *Handler) DeleteArticle(c echo.Context) error {
	err := h.articleUsecase.DeleteArticle(c.Param("articleId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, utils.NewError(err))
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) ListArticle(c echo.Context) error {
	response := h.articleUsecase.ListArticle()
	return c.JSON(http.StatusOK, response)
}
