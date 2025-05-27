package app

import (
	"context"
	"fmt"
	"kiddy-line-processor/config"
	grpclines "kiddy-line-processor/internal/controller/grpc"
	"kiddy-line-processor/internal/controller/http"
	pb "kiddy-line-processor/internal/proto"
	"kiddy-line-processor/internal/repo"
	"kiddy-line-processor/internal/service"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc/reflection"

	"google.golang.org/grpc"
)

type SportsMap = map[string]*service.SportService

// todo: to env
type PullInterval struct {
	Baseball time.Duration
	Football time.Duration
	Soccer   time.Duration
}

type Config struct {
	PullIntervals PullInterval
}

func initLineSportProviders(config Config, sports SportsMap) []*service.LineSportProvider {
	return []*service.LineSportProvider{
		{Sport: sports["baseball"], PullInteval: config.PullIntervals.Baseball},
		{Sport: sports["football"], PullInteval: config.PullIntervals.Football},
		{Sport: sports["soccer"], PullInteval: config.PullIntervals.Soccer},
	}
}

// todo: to ticker
func pullSportLine(ctx context.Context, provider *service.LineSportProvider, wg *sync.WaitGroup) error {
	fmt.Printf("%s start pulling with sleep %s\n", provider.Sport.Sport, provider.PullInteval)
	time.Sleep(provider.PullInteval)
	err := provider.Pull(ctx)

	if err != nil {
		fmt.Println(err)
	}
	if !provider.Synced {
		fmt.Println("Done")
		wg.Done()
	}
	provider.Synced = true
	fmt.Printf("%s pulled!", provider.Sport.Sport)
	return err
}

func runSportPulling(ctx context.Context, provider *service.LineSportProvider, wg *sync.WaitGroup) {
	for {
		pullSportLine(ctx, provider, wg)
	}
}

func runSportsPulling(ctx context.Context, providers []*service.LineSportProvider, wg *sync.WaitGroup) {
	for _, provider := range providers {
		go runSportPulling(ctx, provider, wg)
	}
}

func Run() {
	cfg := config.InitConfig()

	fmt.Println(cfg.GrpcAddr)
	names := []string{
		"baseball",
		"soccer",
		"football",
	}

	sports := make(SportsMap)

	config := Config{
		PullIntervals: PullInterval{
			Baseball: time.Second * 5,
			Soccer:   time.Second * 2,
			Football: time.Second * 1,
		},
	}

	ctx := context.Background()
	redis := repo.Init()

	for _, name := range names {
		sports[name] = service.NewSportService(redis, name)
	}

	providers := initLineSportProviders(config, sports)

	wg := new(sync.WaitGroup)

	wg.Add(len(providers))

	ready := service.NewReadyService(wg)

	runSportsPulling(ctx, providers, wg)

	deps := &service.LineDependencies{
		Sports:       sports,
		ReadyService: ready,
	}

	lineService := &service.LineService{
		Deps: deps,
	}

	httpServer := http.NewServer(":8080", lineService)

	go httpServer.Run()

	fmt.Println("Ждет реади")
	ready.Wait()

	fmt.Println("Иницализация gRPC")
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 8081))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	reflection.Register(grpcServer)
	linesServer := grpclines.NewServer(&service.KiddyLineServiceDeps{
		Sports: sports,
	})
	pb.RegisterSportsLinesServiceServer(grpcServer, linesServer)
	grpcServer.Serve(lis)

}
