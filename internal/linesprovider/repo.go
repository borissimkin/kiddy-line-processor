package linesprovider

import (
	"context"
	"fmt"
	"kiddy-line-processor/internal/storage"
	"strconv"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type LineRepo struct {
	Sport   string
	storage *storage.RedisStorage
}

func NewSportRepo(storage *storage.RedisStorage, sport string) *LineRepo {
	return &LineRepo{
		storage: storage,
		Sport:   sport,
	}
}

func (r *LineRepo) key() string {
	return fmt.Sprintf("lines:%s", r.Sport)
}

func (r *LineRepo) Save(ctx context.Context, coef float64) error {
	_, err := r.storage.RPush(ctx, r.key(), coef).Result()
	if err != nil {
		logrus.Error(err)
	}

	return err
}

func (r *LineRepo) GetLast(ctx context.Context) (CoefItem, error) {
	result := CoefItem{}
	val, err := r.storage.LIndex(ctx, r.key(), -1).Result()
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

func (r *LineRepo) Ready(ctx context.Context) bool {
	if err := r.storage.Ping(ctx).Err(); err != nil {
		logrus.Error(err)
		return false
	}

	return true
}
