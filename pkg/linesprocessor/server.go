package linesprocessor

import (
	"context"
	"errors"
	"fmt"
	"io"
	"kiddy-line-processor/pkg/linesprovider"
	pb "kiddy-line-processor/pkg/proto"
	"math"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type ServerDeps struct {
	Lines linesprovider.LineServiceMap
}

type LinesProcessorServer struct {
	pb.UnimplementedSportsLinesServiceServer

	deps *ServerDeps
	srv  *grpc.Server
}

func NewLinesProcessorServer(deps *ServerDeps) *LinesProcessorServer {
	grpcServer := grpc.NewServer()

	return &LinesProcessorServer{
		deps: deps,
		srv:  grpcServer,
	}
}

func (s *LinesProcessorServer) Run(addr string) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)

		return
	}

	reflection.Register(s.srv)
	pb.RegisterSportsLinesServiceServer(s.srv, s)

	err = s.srv.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}

type PreviousRequest struct {
	Sport []string
}

func round(x float32) float32 {
	const roundPrecision = 100

	return float32(math.Round(float64(x*roundPrecision))) / roundPrecision
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
				log.Error(err)

				return
			}
		case <-ctx.Done():
			return
		}
	}
}

// SubscribeOnSportsLines subscribe to receiving processed sports coefficients.
func (s *LinesProcessorServer) SubscribeOnSportsLines(stream pb.SportsLinesService_SubscribeOnSportsLinesServer) error {
	var (
		prevReq      PreviousRequest
		cancelSender context.CancelFunc
	)

	initialCoef := make(map[string]float32)

	for {
		streamCtx := stream.Context()

		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}

		if err != nil {
			return fmt.Errorf("error receiving previous request: %w", err)
		}

		if cancelSender != nil {
			cancelSender()
		}

		resp := &pb.SubscribeResponse{
			Sports: make(map[string]float32),
		}

		if isSame(prevReq.Sport, req.GetSport()) {
			for _, sport := range req.GetSport() {
				coef, err := s.deps.Lines[sport].GetLast(streamCtx)
				if err != nil {
					return fmt.Errorf("couldn't get coefficient for sport %s: %w", sport, err)
				}

				resp.Sports[sport] = s.getCoefDelta(initialCoef[sport], float32(coef.Coef))
			}
		} else {
			for _, sport := range req.GetSport() {
				coef, err := s.deps.Lines[sport].GetLast(streamCtx)
				if err != nil {
					return fmt.Errorf("error receiving previous request: %w", err)
				}

				rounded := round(float32(coef.Coef))
				resp.Sports[sport] = rounded
				initialCoef[sport] = rounded
			}
		}

		err = stream.Send(resp)
		if err != nil {
			log.Error(err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancelSender = cancel

		prevReq.Sport = req.GetSport()

		go s.SendStream(ctx, stream, req.GetInterval().AsDuration(), req.GetSport(), initialCoef)
	}
}

func (s *LinesProcessorServer) getCoefDelta(a, b float32) float32 {
	return round(a - b)
}
