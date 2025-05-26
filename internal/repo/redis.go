package repo

import (
	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	*redis.Client
}

func Init() *RedisStorage {
	// todo: config
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // No password set
		DB:       0,  // Use default DB
		Protocol: 2,  // Connection protocol
	})

	return &RedisStorage{
		client,
	}
}
