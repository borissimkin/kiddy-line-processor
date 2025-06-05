package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"kiddy-line-processor/config"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type LineSportPuller interface {
	Pull() error
}

type LineSportProvider struct {
	cfg          config.LinesProviderConfig
	SportService *SportService
	PullInteval  time.Duration
	Synced       bool
}

func NewLineSportProvider(
	cfg config.LinesProviderConfig,
	sportService *SportService,
	pullInteval time.Duration,
) *LineSportProvider {
	return &LineSportProvider{
		cfg:          cfg,
		SportService: sportService,
		PullInteval:  pullInteval,
	}
}

func (p *LineSportProvider) Pull(ctx context.Context) error {
	coef, err := p.fetch()
	if err != nil {
		return err
	}

	err = p.SportService.Save(ctx, coef)
	if err != nil {
		return err
	}

	return nil
}

type SportProviderResponse struct {
	Lines map[string]string `json:"lines"`
}

func (p *LineSportProvider) fetch() (float64, error) {
	fmt.Println(p.cfg.Addr())
	resp, err := http.Get(fmt.Sprintf("http://%s/api/v1/lines/%s", p.cfg.Addr(), p.SportService.Sport))

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

	coef := response.Lines[strings.ToUpper(p.SportService.Sport)]

	coefFloat, err := strconv.ParseFloat(coef, 64)

	if err != nil {
		return 0, err
	}

	return coefFloat, nil
}

// todo: вынести
type Line interface {
	Ready(ctx context.Context) bool
}

type LineDependencies struct {
	Sports       map[string]*SportService
	ReadyService *ReadyService
}

type LineService struct {
	Deps *LineDependencies
}

func (s *LineService) Ready(ctx context.Context) bool {
	for _, sport := range s.Deps.Sports {
		if !sport.Ready(ctx) {
			return false
		}
	}

	return s.Deps.ReadyService.IsReady()
}
