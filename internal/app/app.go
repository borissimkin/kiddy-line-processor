package app

import (
	"fmt"
	"kiddy-line-processor/internal/service"
)

func Run() {
	line := service.LineService{}

	resp, err :=line.Fetch()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(resp)
}