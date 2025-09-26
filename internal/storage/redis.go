package storage

import (
	"fmt"
	"kiddy-line-processor/internal/config"

	"github.com/redis/go-redis/v9"
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
