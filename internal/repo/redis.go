package repo

import (
	"fmt"
	"kiddy-line-processor/config"

	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	*redis.Client
}

func Init(cfg config.RedisConfig) *RedisStorage {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password, // No password set
		DB:       cfg.DB,       // Use default DB
		// Protocol: 2,  // Connection protocol
	})

	return &RedisStorage{
		client,
	}
}
