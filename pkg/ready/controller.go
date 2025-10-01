package ready

import (
	"encoding/json"
	"kiddy-line-processor/pkg/config"
	"net/http"

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
	http.HandleFunc("/ready", s.readyHandle)

	log.Fatal(http.ListenAndServe(s.cfg.Addr(), nil))
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
