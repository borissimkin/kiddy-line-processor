package linesprovider

import (
	"context"
)

type CoefItem struct {
	Id   string
	Coef float64
}

type LineRepoInterface interface {
	Save(ctx context.Context, coef float64) error
	GetLast(ctx context.Context) (CoefItem, error)
	Ready(ctx context.Context) bool
}

type LineService struct {
	Sport string
	repo  LineRepoInterface
}

func NewLineService(sport string, repo LineRepoInterface) *LineService {
	return &LineService{
		Sport: sport,
		repo:  repo,
	}
}

func (s *LineService) GetLast(ctx context.Context) (CoefItem, error) {
	return s.repo.GetLast(ctx)
}

func (s *LineService) Save(ctx context.Context, coef float64) error {
	return s.repo.Save(ctx, coef)
}

func (s *LineService) Ready(ctx context.Context) bool {
	return s.repo.Ready(ctx)
}
