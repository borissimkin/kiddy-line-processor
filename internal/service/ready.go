package service

import (
	"sync"
)

// todo: вынести сюда проверку доступности хранилища?
type ReadyService struct {
	Wg    *sync.WaitGroup
	Ready bool
}

func NewReadyService(wg *sync.WaitGroup) *ReadyService {
	return &ReadyService{
		Wg:    wg,
		Ready: false,
	}
}

func (s *ReadyService) Wait() {
	s.Wg.Wait()
	s.Ready = true
}
