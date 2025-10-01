package storage

import (
	"context"
	"fmt"
	"kiddy-line-processor/internal/config"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type RedisStorage struct {
	*redis.Client
}

func Init(cfg config.RedisConfig) *RedisStorage {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	return &RedisStorage{
		client,
	}
}

func (r *RedisStorage) Ready(ctx context.Context) bool {
	err := r.Ping(ctx).Err()
	if err != nil {
		logrus.Error(err)

		return false
	}

	return true
}
