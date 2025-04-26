package repo

type LineStorage interface {
	Save(coef float64) error
	GetAll() ([]CoefItem, error)
	Ready() bool
}
