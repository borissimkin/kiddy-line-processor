package http

import (
	"encoding/json"
	"kiddy-line-processor/internal/service"
	"log"
	"net/http"
)

type Server struct {
	Addr    string
	Service service.Line
}

func NewServer(addr string, service service.Line) *Server {
	return &Server{
		Addr:    addr,
		Service: service,
	}
}

// todo: просто статус нужен
type ReadyResponse struct {
	Ready bool `json:"ready"`
}

func (s *Server) readyHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	isReady := s.Service.Ready(r.Context())

	response := &ReadyResponse{
		Ready: isReady,
	}

	json.NewEncoder(w).Encode(response)
}

func (s *Server) Run() {
	http.HandleFunc("/ready", s.readyHandle)

	log.Fatal(http.ListenAndServe(s.Addr, nil))
}
