package ready

import (
	"context"
	"sync"
	"sync/atomic"
)

type StorageReadyChecker interface {
	Ready(ctx context.Context) bool
}

type LineReadyChecker interface {
	Ready() bool
}

type LinesReadyService struct {
	Wg             *sync.WaitGroup
	ready          atomic.Bool
	Lines          map[string]LineReadyChecker
	storageChecker StorageReadyChecker
}

func (s *LinesReadyService) Ready(ctx context.Context) bool {
	if !s.storageChecker.Ready(ctx) {
		return false
	}

	for _, sport := range s.Lines {
		if !sport.Ready() {
			return false
		}
	}

	return true
}

func (s *LinesReadyService) Wait() {
	s.Wg.Wait()
	s.ready.Store(true)
}

func (s *LinesReadyService) IsReady() bool {
	return s.ready.Load()
}
