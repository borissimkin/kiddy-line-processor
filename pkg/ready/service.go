package ready

import (
	"context"
	"sync"
	"sync/atomic"
)

type storageReadyChecker interface {
	Ready(ctx context.Context) bool
}

// LineSyncedChecker check line is synced.
type LineSyncedChecker interface {
	Synced() bool
}

// LinesReadyService is service for ready http server.
type LinesReadyService struct {
	Wg             *sync.WaitGroup
	ready          atomic.Bool
	Lines          []LineSyncedChecker
	storageChecker storageReadyChecker
}

// NewLinesReadyService constructor.
func NewLinesReadyService(lines []LineSyncedChecker, storageChecker storageReadyChecker) *LinesReadyService {
	return &LinesReadyService{
		Wg:             &sync.WaitGroup{},
		Lines:          lines,
		storageChecker: storageChecker,
		ready:          atomic.Bool{},
	}
}

// Ready checks that storage is ready and lines are synchronized.
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

// Wait waits for lines syncing.
func (s *LinesReadyService) Wait() {
	s.Wg.Wait()
	s.ready.Store(true)
}

// IsReady  .
func (s *LinesReadyService) IsReady() bool {
	return s.ready.Load()
}
