package app

import (
	"fmt"
	grpclines "kiddy-line-processor/internal/controller/grpc"
	"kiddy-line-processor/internal/controller/http"
	pb "kiddy-line-processor/internal/proto"
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
func pullSportLine(provider *service.LineSportProvider, wg *sync.WaitGroup) error {
	fmt.Printf("%s start pulling with sleep %s\n", provider.Sport.Name, provider.PullInteval)
	time.Sleep(provider.PullInteval)
	err := provider.Pull()

	if err != nil {
		fmt.Println(err)
	}
	if !provider.Synced {
		fmt.Println("Done")
		wg.Done()
	}
	provider.Synced = true
	fmt.Printf("%s pulled!", provider.Sport.Name)
	return err
}

func runSportPulling(provider *service.LineSportProvider, wg *sync.WaitGroup) {
	for {
		pullSportLine(provider, wg)
	}
}

func runSportsPulling(providers []*service.LineSportProvider, wg *sync.WaitGroup) {
	for _, provider := range providers {
		go runSportPulling(provider, wg)
	}
}

func Run() {
	config := Config{
		PullIntervals: PullInterval{
			Baseball: time.Second * 1,
			Soccer:   time.Second * 1,
			Football: time.Second * 2,
		},
	}

	names := []string{
		"baseball",
		"soccer",
		"football",
	}

	sports := make(SportsMap)

	for _, name := range names {
		sports[name] = service.NewSportService(name)
	}

	providers := initLineSportProviders(config, sports)

	wg := new(sync.WaitGroup)

	wg.Add(len(providers))

	ready := service.NewReadyService(wg)

	runSportsPulling(providers, wg)

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
