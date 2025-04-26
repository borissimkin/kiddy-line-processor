package service

import (
	"encoding/json"
	"fmt"
	"io"
	"kiddy-line-processor/internal/repo"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type LineSportPuller interface {
	Pull() error
}

type LineSportProvider struct {
	Sport       string
	Storage     repo.LineStorage
	PullInteval time.Duration
	Synced      bool
}

func (p *LineSportProvider) Pull() error {
	coef, err := p.fetch()
	if err != nil {
		return err
	}

	err = p.Storage.Save(coef)
	if err != nil {		
		return err
	}

	p.Synced = true

	return nil
}

type SportProviderResponse struct {
	Lines map[string]string `json:"lines"`
}

func (p *LineSportProvider) fetch() (float64, error) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:8000/api/v1/lines/%s", p.Sport))
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return 0, nil
	}

	var response SportProviderResponse

	err = json.Unmarshal(body, &response)

	if err != nil {
		return 0, err
	}

	coef := response.Lines[strings.ToUpper(p.Sport)]

	coefFloat, err := strconv.ParseFloat(coef, 64)

	if err != nil {
		return 0, err
	}

	return coefFloat, nil
}


type Line interface {
	Ready() bool
}

type LineDependencies struct {
	Providers []*LineSportProvider
}

type LineService struct {
	Deps *LineDependencies
}

func (s *LineService) Ready() bool {
	for _, provider := range s.Deps.Providers {
		if !provider.Storage.Ready() {
			return false
		}

		ready := provider.Synced
		if !ready {
			return false
		}
	}

	return true
}
