package grpc

import (
	"context"
	"fmt"
	"io"
	pb "kiddy-line-processor/internal/proto"
	"kiddy-line-processor/internal/service"
	"time"
)

// todo: добавить валидацию на имя спорта
// todo: здесь нужен только сам сервис, пока напрямую через его зависимости пробуем
type SportsLinesServer struct {
	deps *service.KiddyLineServiceDeps
	pb.UnimplementedSportsLinesServiceServer
}

func NewServer(deps *service.KiddyLineServiceDeps) *SportsLinesServer {
	return &SportsLinesServer{
		deps: deps,
	}
}

type PreviosRequest struct {
	Sport []string
}

func isSame(oldSports []string, sports []string) bool {
	if len(oldSports) != len(sports) {
		return false
	}

	for index, _ := range oldSports {
		if oldSports[index] != sports[index] {
			return false
		}
	}

	return true
}

// {"sport": "soccer", "sport": "football", "interval": "3s"}
// {"sport": "soccer", "sport": "football", "interval": "1s"}
// {"sport": "baseball", "sport": "football", "interval": "5s"}
func (s *SportsLinesServer) SubscribeOnSportsLines(stream pb.SportsLinesService_SubscribeOnSportsLinesServer) error {
	var prevReq PreviosRequest
	initialCoef := make(map[string]float32)
	var cancelSender context.CancelFunc

	for {
		streamCtx := stream.Context()

		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		if cancelSender != nil {
			cancelSender()
		}

		resp := &pb.SubscribeResponse{
			Sports: make(map[string]float32),
		}

		if isSame(prevReq.Sport, req.Sport) {
			for _, sport := range req.Sport {
				coef, err := s.deps.Sports[sport].GetLast(streamCtx)
				if err != nil {
					return err
				}
				resp.Sports[sport] = initialCoef[sport] - float32(coef.Coef)
			}
		} else {
			for _, sport := range req.Sport {
				coef, err := s.deps.Sports[sport].GetLast(streamCtx)
				if err != nil {
					return err
				}
				resp.Sports[sport] = float32(coef.Coef)
				initialCoef[sport] = float32(coef.Coef)
			}
		}

		stream.Send(resp)

		ctx, cancel := context.WithCancel(context.Background())
		cancelSender = cancel

		prevReq.Sport = req.Sport

		go func(req *pb.SubscribeRequest, ctx context.Context) {
			ticker := time.NewTicker(req.Interval.AsDuration())

			for {
				select {
				case <-ctx.Done():
					fmt.Println("отправка остановлена")
					return
				case <-ticker.C:
					resp := &pb.SubscribeResponse{
						Sports: make(map[string]float32),
					}

					for _, sport := range req.Sport {
						coef, _ := s.deps.Sports[sport].GetLast(ctx)

						resp.Sports[sport] = initialCoef[sport] - float32(coef.Coef)
					}

					err := stream.Send(resp)
					if err != nil {
						fmt.Println("ошибка отправки:", err)
						return
					}
				}
			}
		}(req, ctx)
	}
}
