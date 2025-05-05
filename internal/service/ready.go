package service

import (
	"fmt"
	"sync"
)

type ReadyService struct {
	Wg    *sync.WaitGroup
	Ready chan bool
}

func NewReadyService(wg *sync.WaitGroup) *ReadyService {
	return &ReadyService{
		Wg:    wg,
		Ready: make(chan bool),
	}
}

func (s *ReadyService) Wait() {
	go func() {
		s.Wg.Wait()
		fmt.Println("Release wait")
		s.Ready <- true
	}()
}
