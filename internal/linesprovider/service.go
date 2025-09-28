package linesprovider

import (
	"context"
	"sync/atomic"
)

type LineServiceMap = map[string]*LineService

type CoefItem struct {
	Id   string
	Coef float64 // todo: float32?
}

type LineRepoInterface interface {
	Save(ctx context.Context, coef float64) error
	GetLast(ctx context.Context) (CoefItem, error)
}

type LineService struct {
	Sport  string
	synced atomic.Bool
	repo   LineRepoInterface
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

func NewLineServiceMap(sportNames []string, repoFactory func(string) LineRepoInterface) LineServiceMap {
	lines := make(LineServiceMap)
	for _, sport := range sportNames {
		lines[sport] = NewLineService(sport, repoFactory(sport))
	}

	return lines
}

func (s *LineService) SetSynced(val bool) {
	s.synced.Store(val)
}

func (s *LineService) Synced() bool {
	return s.synced.Load()
}
