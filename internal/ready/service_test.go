package ready

import (
	"context"
	readymocks "kiddy-line-processor/internal/ready/mocks"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.uber.org/mock/gomock"
)

func TestReadyService_Ready(t *testing.T) {
	type args struct {
		ctx context.Context
	}

	type MockBehavior func(l []*readymocks.MockLineSyncedChecker, s *readymocks.MockStorageReadyChecker, args args)

	testCases := []struct {
		name         string
		mockBehavior MockBehavior
		args         args
		want         bool
	}{
		{
			name: "should ready",
			mockBehavior: func(l []*readymocks.MockLineSyncedChecker, s *readymocks.MockStorageReadyChecker, args args) {
				s.EXPECT().Ready(args.ctx).Return(true)
				for _, checker := range l {
					checker.EXPECT().Synced().AnyTimes().Return(true)
				}
			},
			args: args{
				ctx: context.Background(),
			},
			want: true,
		},
		{
			name: "should not ready by storage",
			mockBehavior: func(l []*readymocks.MockLineSyncedChecker, s *readymocks.MockStorageReadyChecker, args args) {
				s.EXPECT().Ready(args.ctx).Return(false)
				for _, checker := range l {
					checker.EXPECT().Synced().AnyTimes().Return(true)
				}
			},
			args: args{
				ctx: context.Background(),
			},
			want: false,
		},
		{
			name: "should not ready by one line",
			mockBehavior: func(l []*readymocks.MockLineSyncedChecker, s *readymocks.MockStorageReadyChecker, args args) {
				s.EXPECT().Ready(args.ctx).Return(true)
				for i, checker := range l {
					v := true
					if i == 0 {
						v = false
					}
					checker.EXPECT().Synced().AnyTimes().Return(v)
				}
			},
			args: args{
				ctx: context.Background(),
			},
			want: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLines := make([]*readymocks.MockLineSyncedChecker, 3)
			toPassMockLines := make([]LineSyncedChecker, 3)
			for i := range mockLines {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				checker := readymocks.NewMockLineSyncedChecker(ctrl)
				mockLines[i] = checker
				toPassMockLines[i] = checker
			}

			mockStorageChecker := readymocks.NewMockStorageReadyChecker(ctrl)

			tc.mockBehavior(mockLines, mockStorageChecker, tc.args)
			service := NewLinesReadyService(toPassMockLines, mockStorageChecker)

			got := service.Ready(tc.args.ctx)
			assert.Equal(t, tc.want, got)
		})
	}
}
