package linesprovider

import (
	"context"
	"io"
	"kiddy-line-processor/internal/config"
	pb "kiddy-line-processor/internal/proto"
	"math"
	"net"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type SportsLinesServerDeps struct {
	Lines map[string]*LineService
}

type SportsLinesServer struct {
	deps *SportsLinesServerDeps
	pb.UnimplementedSportsLinesServiceServer
}

func newServer(deps *SportsLinesServerDeps) *SportsLinesServer {
	return &SportsLinesServer{
		deps: deps,
	}
}

func Init(deps *SportsLinesServerDeps, config config.GrpcConfig) error {
	lis, err := net.Listen("tcp", config.Addr())
	if err != nil {
		logrus.Error(err)
		return err
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	linesServer := newServer(deps)
	pb.RegisterSportsLinesServiceServer(grpcServer, linesServer)
	return grpcServer.Serve(lis)
}

type PreviosRequest struct {
	Sport []string
}

func round(x float32) float32 {
	return float32(math.Round(float64(x*100))) / 100
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
				coef, err := s.deps.Lines[sport].GetLast(streamCtx)
				if err != nil {
					return err
				}
				resp.Sports[sport] = round(initialCoef[sport] - float32(coef.Coef))
			}
		} else {
			for _, sport := range req.Sport {
				coef, err := s.deps.Lines[sport].GetLast(streamCtx)
				if err != nil {
					return err
				}

				rounded := round(float32(coef.Coef))
				resp.Sports[sport] = rounded
				initialCoef[sport] = rounded
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
					return
				case <-ticker.C:
					resp := &pb.SubscribeResponse{
						Sports: make(map[string]float32),
					}

					for _, sport := range req.Sport {
						coef, _ := s.deps.Lines[sport].GetLast(ctx)

						resp.Sports[sport] = round(initialCoef[sport] - float32(coef.Coef))
					}

					err := stream.Send(resp)
					if err != nil {
						logrus.Error(err)
						return
					}
				}
			}
		}(req, ctx)
	}
}
