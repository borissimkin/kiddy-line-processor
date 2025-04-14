package service

import (
	"fmt"
	"kiddy-line-processor/internal/repo"
	"log"
	"time"
)

type SportLineCoefFetcher interface {
	Fetch(sport string) (float32, error)
}

type LineSportPuller interface {
	Pull() error
}


type LineSportProvider struct {
	sport        string
	fetcher      SportLineCoefFetcher
	saver        repo.Storage
	PullInterval time.Duration
}

func (p *LineSportProvider) Pull() error {
	log.Println(fmt.Printf("%s pulling...", p.sport))
	coef, err := p.fetcher.Fetch(p.sport)
	if err != nil {
		return err
	}

	err = p.saver.Save(p.sport, coef)
	if err != nil {
		return err
	}
	log.Println(fmt.Printf("%s pulled!", p.sport))
	return nil
}

// type LineResponseWrapper struct {
// 	Lines any `json:"lines"`
// }

// type LineBaseballResponse struct {
// 	Baseball float32 `json:"BASEBALL,string,omitempty"`
// }

// type LineService struct {
// }

// func (s *LineService) Fetch() (*LineBaseballResponse, error) {
// 	resp, err := http.Get("http://localhost:8000/api/v1/lines/baseball")

// 	if err != nil {
// 		return nil, err
// 	}

// 	defer resp.Body.Close()
// 	body, err := io.ReadAll(resp.Body)

// 	if err != nil {
// 		return nil, err
// 	}

// 	payload := &LineBaseballResponse{}

// 	err = json.Unmarshal(body, &LineResponseWrapper{payload})

// 	return payload, err
// }
