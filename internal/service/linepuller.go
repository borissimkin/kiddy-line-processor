package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type LineSportPuller interface {
	Pull() error
}

type LineSportProvider struct {
	Sport       *SportService
	PullInteval time.Duration
	Synced      bool // todo: remove
}

func (p *LineSportProvider) Pull() error {
	coef, err := p.fetch()
	if err != nil {
		return err
	}

	err = p.Sport.Save(coef)
	if err != nil {
		return err
	}

	// p.Synced = true

	return nil
}

type SportProviderResponse struct {
	Lines map[string]string `json:"lines"`
}

func (p *LineSportProvider) fetch() (float64, error) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:8000/api/v1/lines/%s", p.Sport.Name))
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

	coef := response.Lines[strings.ToUpper(p.Sport.Name)]

	coefFloat, err := strconv.ParseFloat(coef, 64)

	if err != nil {
		return 0, err
	}

	return coefFloat, nil
}

// todo: вынести
type Line interface {
	Ready() bool
}

type LineDependencies struct {
	Sports       map[string]*SportService
	ReadyService *ReadyService
}

type LineService struct {
	Deps *LineDependencies
}

func (s *LineService) Ready() bool {
	for _, sport := range s.Deps.Sports {
		if !sport.Ready() {
			return false
		}
	}

	return s.Deps.ReadyService.Ready
}
