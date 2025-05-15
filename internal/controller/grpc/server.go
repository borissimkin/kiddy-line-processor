package grpc

import (
	"fmt"
	"io"
	pb "kiddy-line-processor/internal/proto"
	"kiddy-line-processor/internal/service"
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

func (s *SportsLinesServer) SubscribeOnSportsLines(stream pb.SportsLinesService_SubscribeOnSportsLinesServer) error {
	for {
		req, err := stream.Recv()
		fmt.Println(req)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
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
		}

		stream.Send(resp)

		// for {
		// 	response := &pb.SubscribeResponse{
		// 		Sports: map[string]float32{
		// 			"baseball": 0.3,
		// 		},
		// 	}

		// 	stream.Send(response)

		// 	time.Sleep(time.Second * 8)
		// }

	}
}
