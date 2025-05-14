package grpc

import (
	"fmt"
	"io"
	pb "kiddy-line-processor/internal/proto"
)

type SportsLinesServer struct {
	pb.UnimplementedSportsLinesServiceServer
}

func NewServer() *SportsLinesServer {
	return &SportsLinesServer{}
}

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
		response := &pb.SubscribeResponse{
			Sports: map[string]float32{
				"baseball": 0.3,
			},
		}

		stream.Send(response)
	}
}
