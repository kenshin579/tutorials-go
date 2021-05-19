package redis

import (
	"github.com/kenshin579/tutorials-go/go-redis/blackboard/common/config"
	"github.com/kenshin579/tutorials-go/go-redis/blackboard/domain"
)

type redisBlackBoardStore struct {
	cfg *config.Config
}

func NewRedisBlackBoardStore(cfg *config.Config) domain.BlackBoardStore {
	return &redisBlackBoardStore{
		cfg: cfg,
	}
}
