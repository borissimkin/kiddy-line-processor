package linesprovider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"kiddy-line-processor/internal/config"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type LinesProvider struct {
	cfg          config.LinesProviderConfig
	SportService *LineService
	PullInteval  time.Duration
	Synced       bool
}

func NewLinesProvider(
	cfg config.LinesProviderConfig,
	sportService *LineService,
	pullInteval time.Duration,
) *LinesProvider {
	return &LinesProvider{
		cfg:          cfg,
		SportService: sportService,
		PullInteval:  pullInteval,
	}
}

func (p *LinesProvider) Pull(ctx context.Context) error {
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

type LinesProviderResponse struct {
	Lines map[string]string `json:"lines"`
}

func (p *LinesProvider) fetch() (float64, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s/api/v1/lines/%s", p.cfg.Addr(), p.SportService.Sport))

	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return 0, nil
	}

	var response LinesProviderResponse

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
