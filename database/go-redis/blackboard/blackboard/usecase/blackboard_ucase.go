package usecase

import "github.com/kenshin579/tutorials-go/go-redis/blackboard/domain"

type blackBoardUsecase struct {
	blackBoardStore domain.BlackBoardStore
}

func NewBlackBoardUsecase(ss domain.BlackBoardStore) domain.BlackBoardUsecase {
	return &blackBoardUsecase{
		blackBoardStore: ss,
	}
}
