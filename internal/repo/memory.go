package repo

import (
	"github.com/google/uuid" // todo remove
)

type CoefItem struct {
	Id   string
	Coef float32
}

type MemoryStorage struct {
	Sport string
	coefs []CoefItem
}

func (s *MemoryStorage) Save(coef float32) error {
	s.coefs = append(s.coefs, CoefItem{
		Id:   uuid.New().String(),
		Coef: coef,
	})

	return nil
}

func (s *MemoryStorage) GetAll() ([]CoefItem, error) {
	return s.coefs, nil
}
