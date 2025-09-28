package linesprocessor

import (
	"context"
	"io"
	"kiddy-line-processor/internal/config"
	"kiddy-line-processor/internal/linesprovider"
	pb "kiddy-line-processor/internal/proto"
	"math"
	"net"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type ServerDeps struct {
	Lines linesprovider.LineServiceMap
}

type LinesProcessorServer struct {
	deps *ServerDeps
	pb.UnimplementedSportsLinesServiceServer
}

func newServer(deps *ServerDeps) *LinesProcessorServer {
	return &LinesProcessorServer{
		deps: deps,
	}
}

func Init(deps *ServerDeps, config config.GrpcConfig) error {
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

	for index := range oldSports {
		if oldSports[index] != sports[index] {
			return false
		}
	}

	return true
}

func (s *LinesProcessorServer) getCoefDelta(a, b float32) float32 {
	return round(a - b)
}

func (s *LinesProcessorServer) SendStream(ctx context.Context, stream pb.SportsLinesService_SubscribeOnSportsLinesServer, interval time.Duration, sports []string, initialCoef map[string]float32) {
	ticker := time.NewTicker(interval)

	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			resp := &pb.SubscribeResponse{
				Sports: make(map[string]float32),
			}

			for _, sport := range sports {
				coef, _ := s.deps.Lines[sport].GetLast(ctx)

				resp.Sports[sport] = s.getCoefDelta(initialCoef[sport], float32(coef.Coef))
			}

			err := stream.Send(resp)
			if err != nil {
				logrus.Error(err)
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func (s *LinesProcessorServer) SubscribeOnSportsLines(stream pb.SportsLinesService_SubscribeOnSportsLinesServer) error {
	var prevReq PreviosRequest
	var cancelSender context.CancelFunc
	initialCoef := make(map[string]float32)

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
				resp.Sports[sport] = s.getCoefDelta(initialCoef[sport], float32(coef.Coef))
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

		go s.SendStream(ctx, stream, req.Interval.AsDuration(), req.Sport, initialCoef)
	}
}
