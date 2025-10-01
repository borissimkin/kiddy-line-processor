// Package storage places storage implementations.
package storage

import (
	"context"
	"fmt"
	"kiddy-line-processor/pkg/config"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

// RedisStorage storage.
type RedisStorage struct {
	*redis.Client
}

// Init initialize storage.
func Init(cfg config.RedisConfig) *RedisStorage {
	client := redis.NewClient(&redis.Options{ //nolint:exhaustruct
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	return &RedisStorage{
		client,
	}
}

// Ready check storage connection.
func (r *RedisStorage) Ready(ctx context.Context) bool {
	err := r.Ping(ctx).Err()
	if err != nil {
		logrus.Error(err)

		return false
	}

	return true
}
