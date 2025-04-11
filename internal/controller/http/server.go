package http

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type ServerConfig struct {
	Host string
	Port int
}

type Server struct{}

func readyHandle(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, "Hello from a HandleFunc #1!\n")
}

func (s *Server) Run(cfg ServerConfig) {
	http.HandleFunc("/ready", readyHandle)

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), nil))
}
