package linesprovider

import (
	"context"
	"fmt"
	"sync/atomic"
)

type LineServiceMap = map[string]*LineService

type CoefItem struct {
	Id   string
	Coef float64
}

func NewCoefItem(id string, coef float64) CoefItem {
	return CoefItem{
		Id:   id,
		Coef: coef,
	}
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

// NewLineService constructor.
func NewLineService(sport string, repo LineRepoInterface) *LineService {
	return &LineService{
		Sport:  sport,
		repo:   repo,
		synced: atomic.Bool{},
	}
}

func (s *LineService) GetLast(ctx context.Context) (CoefItem, error) {
	coef, err := s.repo.GetLast(ctx)
	if err != nil {
		return CoefItem{}, fmt.Errorf("failed get last coef: %w", err)
	}

	return coef, nil
}

func (s *LineService) Save(ctx context.Context, coef float64) error {
	err := s.repo.Save(ctx, coef)
	if err != nil {
		return fmt.Errorf("failed to save coef: %w", err)
	}

	return nil
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
