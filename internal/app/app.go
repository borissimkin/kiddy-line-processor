package app

import (
	"fmt"
	"kiddy-line-processor/internal/controller/http"
	"kiddy-line-processor/internal/repo"
	"kiddy-line-processor/internal/service"
	"sync"
	"time"
)

// todo: to env
type PullInterval struct {
	Baseball time.Duration
	Footbal  time.Duration
	Soccer   time.Duration
}

type Config struct {
	PullIntervals PullInterval
}
// todo: кажется нужно сделать вейт груп который в отдельном сервисе проверит что все ченелы в горутинах ready?

func initLineSportProviders(config Config) []*service.LineSportProvider {
	return []*service.LineSportProvider{
		{Sport: "baseball", Storage: &repo.MemoryStorage{Sport: "baseball"}, PullInteval: config.PullIntervals.Baseball},
		{Sport: "football", Storage: &repo.MemoryStorage{Sport: "football"}, PullInteval: config.PullIntervals.Footbal},
		{Sport: "soccer", Storage: &repo.MemoryStorage{Sport: "soccer"}, PullInteval: config.PullIntervals.Soccer},
	}
}

func pullSportLine(provider *service.LineSportProvider) error {
	fmt.Printf("%s start pulling with sleep %s\n", provider.Sport, provider.PullInteval)
	time.Sleep(provider.PullInteval)
	err := provider.Pull()
	fmt.Printf("%s pulled!", provider.Sport)
	return err
}

func runSportPulling(provider *service.LineSportProvider) {
	for {
		pullSportLine(provider)
	}
}

func runSportsPulling(providers []*service.LineSportProvider) {
	for _, provider := range providers {
		go runSportPulling(provider)
	}
}

// func checkLineSynced(providers []service.LineSportProvider) {
// 	for _, provider := range providers {
// 		<-provider.Synced
// 	}

// 	fmt.Println("awdawda ")
// }

func Run() {
	config := Config{
		PullIntervals: PullInterval{
			Baseball: time.Second * 10,
			Soccer:   time.Second * 3,
			Footbal:  time.Second * 5,
		},
	}

	providers := initLineSportProviders(config)

	var wg sync.WaitGroup

	runSportsPulling(providers)

	deps := &service.LineDependencies{
		Providers: providers,
	}

	lineService := &service.LineService{
		Deps: deps,
	}

	httpServer := http.NewServer(":8080", lineService)

	httpServer.Run()
}
