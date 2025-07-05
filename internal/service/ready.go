package service

import (
	"sync"
	"sync/atomic"
)

type ReadyService struct {
	Wg    *sync.WaitGroup
	Ready atomic.Bool
}

func NewReadyService(wg *sync.WaitGroup) *ReadyService {
	return &ReadyService{
		Wg: wg,
	}
}

func (s *ReadyService) Wait() {
	s.Wg.Wait()
	s.Ready.Store(true)
}

func (s *ReadyService) IsReady() bool {
	return s.Ready.Load()
}
