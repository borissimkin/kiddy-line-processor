package ready

import (
	"context"
	"sync"
	"sync/atomic"
)

type StorageReadyChecker interface {
	Ready(ctx context.Context) bool
}

type LineSyncedChecker interface {
	Synced() bool
}

type LinesReadyService struct {
	Wg             *sync.WaitGroup
	ready          atomic.Bool
	Lines          []LineSyncedChecker
	storageChecker StorageReadyChecker
}

func NewLinesReadyService(lines []LineSyncedChecker, storageChecker StorageReadyChecker) *LinesReadyService {
	return &LinesReadyService{
		Wg:             &sync.WaitGroup{},
		Lines:          lines,
		storageChecker: storageChecker,
	}
}

func (s *LinesReadyService) Ready(ctx context.Context) bool {
	if !s.storageChecker.Ready(ctx) {
		return false
	}

	for _, line := range s.Lines {
		if !line.Synced() {
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
