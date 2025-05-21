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

// func (s *SportsLinesServer) runSendDeltas(req *pb.SubscribeRequest) {

// }

// func (s *SportsLinesServer) sendLine(stream pb.SportsLinesService_SubscribeOnSportsLinesServer, ch <-chan pb.SubscribeRequest) error {
// 	for {

// 	}
// }

type PreviosRequest struct {
	Sport    []string
	Interval time.Duration
}

func (s *SportsLinesServer) SubscribeOnSportsLines(stream pb.SportsLinesService_SubscribeOnSportsLinesServer) error {
	var prevReq PreviosRequest
	initialCoef := make(map[string]float32)
	var cancelSender context.CancelFunc

	for {
		// reqCh := make(chan pb.SubscribeRequest)
		stream.Context()
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

		for _, sport := range req.Sport {
			coef, err := s.deps.Sports[sport].GetLast()
			if err != nil {
				return err
			}
			resp.Sports[sport] = float32(coef.Coef)
			initialCoef[sport] = float32(coef.Coef)
		}

		stream.Send(resp)

		ctx, cancel := context.WithCancel(context.Background())
		cancelSender = cancel

		prevReq.Interval = req.Interval.AsDuration()
		prevReq.Sport = req.Sport

		go func(req *pb.SubscribeRequest, ctx context.Context) {
			ticker := time.NewTicker(time.Second * 5)

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
						coef, _ := s.deps.Sports[sport].GetLast()

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
