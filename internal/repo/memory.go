package repo

import (
	"github.com/google/uuid"
)

type CoefItem struct {
	Id   string
	Coef float64
}

type MemoryStorage struct {
	Sport string
	coefs []CoefItem
}

func (s *MemoryStorage) Save(coef float64) error {
	s.coefs = append(s.coefs, CoefItem{
		Id:   uuid.New().String(),
		Coef: coef,
	})

	return nil
}

func (s *MemoryStorage) GetAll() ([]CoefItem, error) {
	return s.coefs, nil
}

func (s *MemoryStorage) Ready() bool {
	return true
}
