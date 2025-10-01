// Package app run and initialize main service.
package app

import (
	"context"
	cfg "kiddy-line-processor/pkg/config"
	"kiddy-line-processor/pkg/linesprocessor"
	"kiddy-line-processor/pkg/linesprovider"
	"kiddy-line-processor/pkg/ready"
	"kiddy-line-processor/pkg/storage"

	log "github.com/sirupsen/logrus"
)

// Run init and run app.
func Run() {
	config := cfg.InitConfig()

	SetLogger(config.Logger.Level)

	sportNames := []string{
		"baseball",
		"soccer",
		"football",
	}

	redis := storage.Init(config.Redis)

	repoFactory := func(sport string) linesprovider.LineRepoInterface {
		return linesprovider.NewSportRepo(redis, sport)
	}

	lineServiceMap := linesprovider.NewLineServiceMap(sportNames, repoFactory)

	lineSyncedCheckers := make([]ready.LineSyncedChecker, 0)
	for _, v := range lineServiceMap {
		lineSyncedCheckers = append(lineSyncedCheckers, v)
	}

	readyService := ready.NewLinesReadyService(lineSyncedCheckers, redis)
	readyService.Wg.Add(len(lineSyncedCheckers))

	ctx := context.Background()

	httpServer := ready.NewServer(config.HTTP, readyService)
	go httpServer.Run()

	linesPullService := linesprovider.InitLinesPullService(config, lineServiceMap)
	linesPullService.StartPulling(ctx, readyService.Wg)

	log.Info("wait for lines syncing...")
	readyService.Wait()
	log.Info("gRPC initializing...")

	deps := &linesprocessor.ServerDeps{
		Lines: lineServiceMap,
	}

	linesProcessorSrv := linesprocessor.NewLinesProcessorServer(deps)

	log.Info("gRPC run...")
	linesProcessorSrv.Run(ctx, config.Grpc.Addr())
}
