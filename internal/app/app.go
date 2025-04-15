package app

import (
	"fmt"
	"kiddy-line-processor/internal/controller/http"
	"kiddy-line-processor/internal/repo"
	"kiddy-line-processor/internal/service"
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

func initLineSportProviders(config Config) []service.LineSportProvider {
	return []service.LineSportProvider{
		{Sport: "baseball", Storage: &repo.MemoryStorage{Sport: "baseball"}, PullInteval: config.PullIntervals.Baseball},
		{Sport: "football", Storage: &repo.MemoryStorage{Sport: "football"}, PullInteval: config.PullIntervals.Footbal},
		{Sport: "soccer", Storage: &repo.MemoryStorage{Sport: "soccer"}, PullInteval: config.PullIntervals.Soccer},
	}
}

func pullSportLine(provider service.LineSportProvider) {
	fmt.Println(fmt.Printf("%s start pulling with sleep %s", provider.Sport, provider.PullInteval))
	time.Sleep(provider.PullInteval)
	provider.Pull()
	fmt.Println(fmt.Printf("%s pulled!", provider.Sport))
}

func runSportPulling(provider service.LineSportProvider) {
	for {
		pullSportLine(provider)
	}
}

func runSportsPulling(providers []service.LineSportProvider) {
	for _, provider := range providers {
		runSportPulling(provider)
	}

}

func Run() {
	config := Config{
		PullIntervals: PullInterval{
			Baseball: time.Second * 3,
			Soccer:   time.Second * 6,
			Footbal:  time.Second * 8,
		},
	}

	providers := initLineSportProviders(config)

	runSportsPulling(providers)

	httpServer := http.NewServer(":8080")

	httpServer.Run()
}
