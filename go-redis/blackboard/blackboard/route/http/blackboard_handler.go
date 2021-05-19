package http

import (
	"net/http"

	"github.com/kenshin579/tutorials-go/go-redis/blackboard/domain"
	"github.com/labstack/echo"
)

type blackBoardHandler struct {
	blackBoardUsecase domain.BlackBoardUsecase
}

func NewBlackBoardHandler(g *echo.Group, su domain.BlackBoardUsecase) *blackBoardHandler {
	handler := &blackBoardHandler{
		blackBoardUsecase: su,
	}

	blackboard := g.Group("/blackboard")

	blackboard.POST("/job", handler.Set)

	return handler
}

func (s *blackBoardHandler) Set(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNoContent)
}
