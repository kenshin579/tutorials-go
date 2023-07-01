package handler

import "github.com/labstack/echo"

func (h *Handler) Register(v1 *echo.Group) {
	articles := v1.Group("/articles")
	articles.POST("", h.CreateArticle)
	articles.DELETE("/:articleId", h.DeleteArticle)
	articles.GET("", h.ListArticle)
	articles.GET("/:articleId", h.GetArticle)
}
