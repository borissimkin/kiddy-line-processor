package linesprovider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"kiddy-line-processor/pkg/config"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type LinesProvider struct {
	cfg          config.LinesProviderConfig
	lineService  *LineService
	pullInterval time.Duration
	client       *http.Client
}

type LinesProviders = map[string]*LinesProvider

func NewLinesProvider(
	cfg config.LinesProviderConfig,
	lineService *LineService,
	interval time.Duration,
) *LinesProvider {
	const timeout = time.Second * 10

	return &LinesProvider{
		cfg:          cfg,
		lineService:  lineService,
		pullInterval: interval,
		client: &http.Client{ //nolint:exhaustruct
			Timeout: timeout,
		},
	}
}

func (p *LinesProvider) Pull(ctx context.Context) error {
	coef, err := p.fetch(ctx)
	if err != nil {
		return err
	}

	err = p.lineService.Save(ctx, coef)
	if err != nil {
		return err
	}

	return nil
}

func (p *LinesProvider) StartPulling(ctx context.Context, waitGroup *sync.WaitGroup) {
	ctxLogger := log.WithFields(log.Fields{
		"provider": p.lineService.Sport,
		"interval": p.pullInterval.String(),
	})

	ctxLogger.Info("Start pulling")

	ticker := time.NewTicker(p.pullInterval)

	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := p.Pull(ctx)
			if err != nil {
				ctxLogger.Error(err)

				if !p.lineService.Synced() {
					return
				}
			}

			ctxLogger.Info("Pulled successfully")

			if !p.lineService.Synced() {
				waitGroup.Done()
				p.lineService.SetSynced(true)
				ctxLogger.Info("Is synced")
			}

		case <-ctx.Done():
			ctxLogger.Info("Stop pulling")
		}
	}
}

type LinesProviderResponse struct {
	Lines map[string]string `json:"lines"`
}

func (p *LinesProvider) fetch(ctx context.Context) (float64, error) {
	url := fmt.Sprintf("http://%s/api/v1/lines/%s", p.cfg.Addr(), p.lineService.Sport)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed http get in fetch: %w", err)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Errorf("Failed to close response body: %v", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil
	}

	var response LinesProviderResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, fmt.Errorf("faied unmarshal in fetch: %w", err)
	}

	coef := response.Lines[strings.ToUpper(p.lineService.Sport)]

	coefFloat, err := strconv.ParseFloat(coef, 64)
	if err != nil {
		return 0, fmt.Errorf("failed ParseFloat in fetch: %w", err)
	}

	return coefFloat, nil
}

type LinesPullService struct {
	linesProviders []*LinesProvider
}

func InitLinesPullService(config config.Config, lines LineServiceMap) *LinesPullService {
	return &LinesPullService{
		linesProviders: []*LinesProvider{
			NewLinesProvider(config.LinesProvider, lines["baseball"], config.PullInterval.Baseball),
			NewLinesProvider(config.LinesProvider, lines["soccer"], config.PullInterval.Soccer),
			NewLinesProvider(config.LinesProvider, lines["football"], config.PullInterval.Football),
		},
	}
}

func (s *LinesPullService) StartPulling(ctx context.Context, wg *sync.WaitGroup) {
	for _, provider := range s.linesProviders {
		go provider.StartPulling(ctx, wg)
	}
}
