package ready_test

import (
	"context"
	"kiddy-line-processor/pkg/ready"
	readymocks "kiddy-line-processor/pkg/ready/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.uber.org/mock/gomock"
)

func TestReadyService_Ready(t *testing.T) {
	t.Parallel()

	type MockBehavior func(ctx context.Context, l []*readymocks.MockLineSyncedChecker, s *readymocks.MockStorageReadyChecker)

	testCases := []struct {
		name         string
		mockBehavior MockBehavior
		want         bool
	}{
		{
			name: "should ready",
			mockBehavior: func(ctx context.Context, l []*readymocks.MockLineSyncedChecker, s *readymocks.MockStorageReadyChecker) {
				s.EXPECT().Ready(ctx).Return(true)

				for _, checker := range l {
					checker.EXPECT().Synced().AnyTimes().Return(true)
				}
			},
			want: true,
		},
		{
			name: "should not ready by storage",
			mockBehavior: func(ctx context.Context, l []*readymocks.MockLineSyncedChecker, s *readymocks.MockStorageReadyChecker) {
				s.EXPECT().Ready(ctx).Return(false)

				for _, checker := range l {
					checker.EXPECT().Synced().AnyTimes().Return(true)
				}
			},
			want: false,
		},
		{
			name: "should not ready by one line",
			mockBehavior: func(ctx context.Context, l []*readymocks.MockLineSyncedChecker, s *readymocks.MockStorageReadyChecker) {
				s.EXPECT().Ready(ctx).Return(true)

				for i, checker := range l {
					v := true
					if i == 0 {
						v = false
					}

					checker.EXPECT().Synced().AnyTimes().Return(v)
				}
			},
			want: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			countLines := 3
			mockLines := make([]*readymocks.MockLineSyncedChecker, countLines)
			lineCheckers := make([]ready.LineSyncedChecker, countLines)

			for index := range mockLines {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				checker := readymocks.NewMockLineSyncedChecker(ctrl)
				mockLines[index] = checker
				lineCheckers[index] = checker
			}

			mockStorageChecker := readymocks.NewMockStorageReadyChecker(ctrl)

			ctx := context.Background()
			testCase.mockBehavior(ctx, mockLines, mockStorageChecker)
			service := ready.NewLinesReadyService(lineCheckers, mockStorageChecker)

			got := service.Ready(ctx)
			assert.Equal(t, testCase.want, got)
		})
	}
}

func TestReadyService_Wait(t *testing.T) {
	t.Parallel()

	type MockBehavior func(ctx context.Context, l []*readymocks.MockLineSyncedChecker, s *readymocks.MockStorageReadyChecker)

	testCases := []struct {
		name         string
		mockBehavior MockBehavior
		want         bool
	}{
		{
			name: "should wait and set ready true",
			mockBehavior: func(ctx context.Context, l []*readymocks.MockLineSyncedChecker, s *readymocks.MockStorageReadyChecker) {
				s.EXPECT().Ready(ctx).AnyTimes().Return(true)

				for _, checker := range l {
					checker.EXPECT().Synced().AnyTimes().Return(true)
				}
			},
			want: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			countLines := 3
			mockLines := make([]*readymocks.MockLineSyncedChecker, countLines)
			lineCheckers := make([]ready.LineSyncedChecker, countLines)

			for index := range mockLines {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				checker := readymocks.NewMockLineSyncedChecker(ctrl)
				mockLines[index] = checker
				lineCheckers[index] = checker
			}

			mockStorageChecker := readymocks.NewMockStorageReadyChecker(ctrl)

			ctx := context.Background()
			testCase.mockBehavior(ctx, mockLines, mockStorageChecker)
			service := ready.NewLinesReadyService(lineCheckers, mockStorageChecker)

			service.Wg.Add(countLines)

			now := time.Now()

			go func() {
				time.Sleep(10 * time.Millisecond)

				for range countLines {
					service.Wg.Done()
				}
			}()

			service.Wait()

			after := time.Now()

			assert.GreaterOrEqual(t, after.Sub(now), 10*time.Millisecond)
			assert.Equal(t, testCase.want, service.IsReady())
		})
	}
}
