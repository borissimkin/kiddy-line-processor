package repo

import "context"

type LineStorage interface {
	Save(ctx context.Context, coef float64) error
	// GetAll() ([]CoefItem, error)
	GetLast(ctx context.Context) (CoefItem, error)
	Ready(ctx context.Context) bool
}
