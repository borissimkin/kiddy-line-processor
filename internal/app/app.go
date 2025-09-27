package app

import (
	"context"
	cfg "kiddy-line-processor/internal/config"
	"kiddy-line-processor/internal/linesprocessor"
	"kiddy-line-processor/internal/linesprovider"
	"kiddy-line-processor/internal/ready"
	"kiddy-line-processor/internal/storage"

	log "github.com/sirupsen/logrus"
)

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

	httpServer := ready.NewServer(config.Http, readyService)
	go httpServer.Run()

	ctx, _ := context.WithCancel(context.Background()) // todo: cancel
	linesPullService := linesprovider.InitLinesPullService(config, lineServiceMap)
	linesPullService.StartPulling(ctx, readyService.Wg)

	log.Info("wait for lines syncing...")
	readyService.Wait()
	log.Info("gRPC initializing...")

	deps := &linesprocessor.ServerDeps{
		Lines: lineServiceMap,
	}
	err := linesprocessor.Init(deps, config.Grpc)

	// todo: graceful shutdown
	if err != nil {
		log.Fatal(err)
	}
}
