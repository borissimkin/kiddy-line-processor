package app

import (
	"context"
	"fmt"
	cfg "kiddy-line-processor/config"
	grpclines "kiddy-line-processor/internal/controller/grpc"
	"kiddy-line-processor/internal/controller/http"
	"kiddy-line-processor/internal/repo"
	"kiddy-line-processor/internal/service"
	"log"
	"sync"
	"time"
)

type SportsMap = map[string]*service.SportService

func initLineSportProviders(config cfg.Config, sports SportsMap) []*service.LineSportProvider {
	return []*service.LineSportProvider{
		service.NewLineSportProvider(config.LinesProvider, sports["baseball"], config.PullInterval.Baseball),
		service.NewLineSportProvider(config.LinesProvider, sports["soccer"], config.PullInterval.Soccer),
		service.NewLineSportProvider(config.LinesProvider, sports["football"], config.PullInterval.Football),
	}
}

// todo: to ticker
func pullSportLine(ctx context.Context, provider *service.LineSportProvider, wg *sync.WaitGroup) error {
	fmt.Printf("%s start pulling with sleep %s\n", provider.SportService.Sport, provider.PullInteval)
	time.Sleep(provider.PullInteval)
	err := provider.Pull(ctx)

	if err != nil {
		fmt.Println(err)
		return err
	}
	if !provider.Synced {
		fmt.Println("Done")
		wg.Done()
	}
	provider.Synced = true
	fmt.Printf("%s pulled!", provider.SportService.Sport)
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
	config := cfg.InitConfig()

	names := []string{
		"baseball",
		"soccer",
		"football",
	}

	sports := make(SportsMap)

	ctx := context.Background()
	redis := repo.Init(config.RedisConfig)

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

	httpServer := http.NewServer(config.Http, lineService)

	go httpServer.Run()

	fmt.Println("Ждет реади")
	ready.Wait()

	fmt.Println("Иницализация gRPC")
	err := grpclines.Init(&service.KiddyLineServiceDeps{
		Sports: sports,
	}, config.Grpc)

	if err != nil {
		log.Fatal(err)
	}
}
