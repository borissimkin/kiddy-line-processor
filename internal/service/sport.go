package service

import (
	"context"
	"kiddy-line-processor/internal/repo"
)

type SportService struct {
	Sport   string
	storage repo.LineStorage
}

func NewSportService(redis *repo.RedisStorage, sport string) *SportService {
	return &SportService{
		Sport:   sport,
		storage: repo.NewSportRepo(redis, sport),
	}
}

func (s *SportService) GetLast(ctx context.Context) (repo.CoefItem, error) {
	return s.storage.GetLast(ctx)
}

func (s *SportService) Save(ctx context.Context, coef float64) error {
	return s.storage.Save(ctx, coef)
}

func (s *SportService) Ready(ctx context.Context) bool {
	return s.storage.Ready(ctx)
}
