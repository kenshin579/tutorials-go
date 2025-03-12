package handler

import (
	"net/http"

	"github.com/kenshin579/tutorials-go/go-validation/article/exception"

	"github.com/kenshin579/tutorials-go/go-validation/article/utils"

	"github.com/kenshin579/tutorials-go/go-validation/article/model"
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
		return exception.ErrBinding
	}

	if err := c.Validate(request); err != nil {
		return err
	}

	response, err := h.articleUsecase.CreateArticle(request)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, response)
}

func (h *Handler) GetArticle(c echo.Context) error {
	response, err := h.articleUsecase.GetArticle(c.Param("articleId"))
	if err != nil {
		return err
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
