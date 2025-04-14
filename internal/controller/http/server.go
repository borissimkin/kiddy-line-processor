package http

import (
	"io"
	"log"
	"net/http"
)

type Server struct {
	Addr string
}

func NewServer(Addr string) *Server {
	return &Server{
		Addr: Addr,
	}
}

func readyHandle(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, "Hello from a HandleFunc #1!\n")
}
	
func (s *Server) Run() {
	http.HandleFunc("/ready", readyHandle)

	log.Fatal(http.ListenAndServe(s.Addr, nil))
}
