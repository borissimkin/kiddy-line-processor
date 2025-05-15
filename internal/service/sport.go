package service

import (
	"kiddy-line-processor/internal/repo"
)

type SportService struct {
	Name    string
	storage repo.LineStorage
}

func NewSportService(name string) *SportService {
	return &SportService{
		Name:    name,
		storage: &repo.MemoryStorage{Sport: name},
	}
}

func (s *SportService) GetLast() (repo.CoefItem, error) {
	return s.storage.GetLast()
}

func (s *SportService) Save(coef float64) error {
	return s.storage.Save(coef)
}

func (s *SportService) Ready() bool {
	return s.storage.Ready()
}
