package repo

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type CoefItem struct {
	Id   string
	Coef float64
}

type SportRepo struct {
	Sport  string
	client *RedisStorage
}

func NewSportRepo(client *RedisStorage, sport string) *SportRepo {
	return &SportRepo{
		client: client,
		Sport:  sport,
	}
}

func (r *SportRepo) key() string {
	return fmt.Sprintf("lines:%s", r.Sport)
}

func (r *SportRepo) Save(ctx context.Context, coef float64) error {
	_, err := r.client.RPush(ctx, r.key(), coef).Result()
	if err != nil {
		logrus.Error(err)
	}

	return err
}

func (r *SportRepo) GetLast(ctx context.Context) (CoefItem, error) {
	result := CoefItem{}
	val, err := r.client.LIndex(ctx, r.key(), -1).Result()
	if err != nil {
		logrus.Error(err)
		return result, err
	}

	result.Id = uuid.New().String()
	res, err := strconv.ParseFloat(val, 32)
	if err != nil {
		return result, err
	}

	result.Coef = res
	return result, nil
}

func (r *SportRepo) Ready(ctx context.Context) bool {
	if err := r.client.Ping(ctx).Err(); err != nil {
		logrus.Error(err)
		return false
	}

	return true
}
