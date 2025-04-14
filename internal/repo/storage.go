package repo

type LineStorage interface {
	Save(key string, coef float32) error
	GetAll() ([]float32, error)
}
