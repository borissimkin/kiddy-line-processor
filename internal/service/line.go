package service

import (
	"kiddy-line-processor/internal/repo"
	"time"
)

type LineSportPuller interface {
	Pull() error
}

type LineSportProvider struct {
	Sport       string
	Storage     repo.LineStorage
	PullInteval time.Duration
}

func (p *LineSportProvider) Pull() error {
	// log.Println(fmt.Printf("%s pulling...", p.Sport))
	coef, err := p.fetch()
	if err != nil {
		return err
	}

	err = p.Storage.Save(coef)
	if err != nil {
		return err
	}
	// log.Println(fmt.Printf("%s pulled!", p.Sport))
	return nil
}

func (p *LineSportProvider) fetch() (float32, error) {
	return 1.2, nil
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
