package ready

import (
	"encoding/json"
	"errors"
	"kiddy-line-processor/pkg/config"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type Server struct {
	cfg     config.HttpConfig
	service *LinesReadyService
}

func NewServer(cfg config.HttpConfig, service *LinesReadyService) *Server {
	return &Server{
		cfg:     cfg,
		service: service,
	}
}

type ReadyResponse struct {
	Ready bool `json:"ready"`
}

func (s *Server) Run() {
	const (
		idleTimeout  = 120 * time.Second
		writeTimeout = 10 * time.Second
		readTimeout  = 5 * time.Second
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/ready", s.readyHandle)

	srv := &http.Server{ //nolint:exhaustruct
		Addr:         s.cfg.Addr(),
		Handler:      mux,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	err := srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrAbortHandler) {
		log.Fatalf("HTTP server failed: %v", err)
	}
}

func (s *Server) readyHandle(writer http.ResponseWriter, r *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	isReady := s.service.Ready(r.Context())

	response := &ReadyResponse{
		Ready: isReady,
	}

	err := json.NewEncoder(writer).Encode(response)
	if err != nil {
		log.Error(err)
	}
}
