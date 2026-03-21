package server

import (
	"github.com/MdRasB/LogLine/internal/handler"
)

func (s *Server) registerRoutes() {
	s.mux.HandleFunc("/health", handler.HandleHealth)
	s.mux.HandleFunc("/ingest", handler.HandleIngest)
}