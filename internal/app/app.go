package app

import (
	"fmt"
	"kiddy-line-processor/internal/controller/http"
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

func InitLineSportProviders() []service.LineSportProvider {
	
}

func Run() {
	config := Config{
		PullIntervals: PullInterval{
			Baseball: time.Second * 3,
			Soccer:   time.Second * 6,
			Footbal:  time.Second * 8,
		},
	}

	httpServer := http.NewServer(":8080")

	go httpServer.Run()

	line := service.LineService{}

	resp, err := line.Fetch()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(resp)
}
