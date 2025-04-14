package repo

type MemoryStorage struct {
	coefs []float32
}

func (s *MemoryStorage) Save(coef float32) error {
	s.coefs = append(s.coefs, coef)
	return nil
}

func (s *MemoryStorage) GetAll() ([]float32, error) {
	return s.coefs, nil
}
