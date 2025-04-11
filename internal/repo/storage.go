package repo

type Storage interface {
	Save(key string, coef float32) error
	Get(key string) (float32, error)
}

// type Saver interface {
// Save(key string, coef float32) error
// }
