package repo

type LineStorage interface {
	Save(coef float64) error
	GetAll() ([]CoefItem, error)
	GetLast() (CoefItem, error)
	Ready() bool
}
