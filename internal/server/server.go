package server

import (
	"net/http"
)

type Server struct {
	addr string
	mux  *http.ServeMux
}

func NewServer(addr string) *Server {
	mux := http.NewServeMux()

	s := &Server{
		addr: addr,
		mux:  mux,
	}
	
	s.registerRoutes()

	return s
}

func (s *Server) Start() error {
	return http.ListenAndServe(s.addr, s.mux)
}