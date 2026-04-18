package server

import (
	"github.com/MdRasB/LogLine/internal/handler"
)

func (s *Server) registerRoutes() {
	ingestHandler := handler.NewIngestHandler(s.logStore)
	s.mux.HandleFunc("/health", handler.HandleHealth)
	s.mux.HandleFunc("/ingest", ingestHandler.Handle)
}
