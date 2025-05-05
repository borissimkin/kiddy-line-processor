package app

import (
	"fmt"
	"kiddy-line-processor/internal/controller/http"
	"kiddy-line-processor/internal/repo"
	"kiddy-line-processor/internal/service"
	"runtime"
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

func pullSportLine(provider *service.LineSportProvider, wg *sync.WaitGroup) error {
	fmt.Printf("%s start pulling with sleep %s\n", provider.Sport, provider.PullInteval)
	time.Sleep(provider.PullInteval)
	err := provider.Pull()
	// todo: check err
	if !provider.Synced {
		fmt.Println("Done")
		wg.Done()
	}
	provider.Synced = true
	fmt.Printf("%s pulled!", provider.Sport)
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

	wg := new(sync.WaitGroup)

	wg.Add(len(providers))

	ready := service.NewReadyService(wg)

	runSportsPulling(providers, wg)

	ready.Wait()

	deps := &service.LineDependencies{
		Providers:    providers,
		ReadyService: ready,
	}

	lineService := &service.LineService{
		Deps: deps,
	}

	httpServer := http.NewServer(":8080", lineService)

	go httpServer.Run()

	fmt.Println("Ждет реади")
	awd := <-ready.Ready
	fmt.Println("Иницализация gRPC")
	fmt.Println(awd)
	runtime.Goexit()

}
