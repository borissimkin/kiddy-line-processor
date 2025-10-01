// Package ready is a package for checking that service is ready for processed coefficients streaming It provides
// http server with ready status information.
package ready

import (
	"encoding/json"
	"errors"
	"kiddy-line-processor/pkg/config"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// Server defines configuration for http server.
type Server struct {
	cfg     config.HTTPConfig
	service *LinesReadyService
}

// NewServer constructor.
func NewServer(cfg config.HTTPConfig, service *LinesReadyService) *Server {
	return &Server{
		cfg:     cfg,
		service: service,
	}
}

type response struct {
	Ready bool `json:"ready"`
}

// Run starts server.
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

	resp := &response{
		Ready: isReady,
	}

	err := json.NewEncoder(writer).Encode(resp)
	if err != nil {
		log.Error(err)
	}
}
