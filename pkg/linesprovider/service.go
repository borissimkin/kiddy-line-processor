package linesprovider

import (
	"context"
	"fmt"
	"sync/atomic"
)

// LineServiceMap is a map of sport and its service.
type LineServiceMap = map[string]*LineService

// CoefItem is an object of saved sport coefficient.
type CoefItem struct {
	ID   string
	Coef float64
}

// NewCoefItem constructor.
func NewCoefItem(id string, coef float64) CoefItem {
	return CoefItem{
		ID:   id,
		Coef: coef,
	}
}

// LineRepoInterface is an interface for repository.
type LineRepoInterface interface {
	Save(ctx context.Context, coef float64) error
	GetLast(ctx context.Context) (CoefItem, error)
}

// LineService is a service for sport line.
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

// GetLast returns last saved coefficient from repository.
func (s *LineService) GetLast(ctx context.Context) (CoefItem, error) {
	coef, err := s.repo.GetLast(ctx)
	if err != nil {
		return CoefItem{}, fmt.Errorf("failed get last coef: %w", err)
	}

	return coef, nil
}

// Save saves coefficient to repository.
func (s *LineService) Save(ctx context.Context, coef float64) error {
	err := s.repo.Save(ctx, coef)
	if err != nil {
		return fmt.Errorf("failed to save coef: %w", err)
	}

	return nil
}

// NewLineServiceMap constructor.
func NewLineServiceMap(sportNames []string, repoFactory func(string) LineRepoInterface) LineServiceMap {
	lines := make(LineServiceMap)
	for _, sport := range sportNames {
		lines[sport] = NewLineService(sport, repoFactory(sport))
	}

	return lines
}

// SetSynced sets synced.
func (s *LineService) SetSynced(val bool) {
	s.synced.Store(val)
}

// Synced checks line was synced.
func (s *LineService) Synced() bool {
	return s.synced.Load()
}
