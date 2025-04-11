package app

import (
	"fmt"
	"kiddy-line-processor/internal/controller/http"
	"kiddy-line-processor/internal/service"
)

func Run() {
	// todo: config

	httpServer := http.NewServer(":8080")

	go httpServer.Run()

	line := service.LineService{}

	resp, err := line.Fetch()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(resp)
}
