package repo

type LineStorage interface {
	Save(coef float32) error
	GetAll() ([]CoefItem, error)
}
