package linesprovider

import (
	"context"
	"fmt"
	"kiddy-line-processor/pkg/storage"
	"strconv"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
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

func (r *LineRepo) Save(ctx context.Context, coef float64) error {
	_, err := r.storage.RPush(ctx, r.key(), coef).Result()
	if err != nil {
		e := fmt.Errorf("failed RPush in Save: %w", err)
		log.Error(e)

		return e
	}

	return nil
}

func (r *LineRepo) GetLast(ctx context.Context) (CoefItem, error) {
	val, err := r.storage.LIndex(ctx, r.key(), -1).Result()
	if err != nil {
		e := fmt.Errorf("failed LIndex in GetLast: %w", err)
		log.Error(e)

		return CoefItem{}, e
	}

	res, err := strconv.ParseFloat(val, 32)
	if err != nil {
		return CoefItem{}, fmt.Errorf("failed ParseFloat in GetLast: %w", err)
	}

	result := NewCoefItem(uuid.New().String(), res)

	return result, nil
}

func (r *LineRepo) Ready(ctx context.Context) bool {
	err := r.storage.Ping(ctx).Err()
	if err != nil {
		log.Error(err)

		return false
	}

	return true
}

func (r *LineRepo) key() string {
	return "lines:" + r.Sport
}
