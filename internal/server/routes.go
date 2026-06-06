package server

import (
	"github.com/MdRasB/LogLine/internal/auth"
	"github.com/MdRasB/LogLine/internal/handler"
)

func (s *Server) registerRoutes() {

	authService := auth.NewService(
		&s.userStore,
		&s.sessionStore,
	)

	ingestHandler := handler.NewIngestHandler(&s.logStore)
	logHandler := handler.NewLogHandler(&s.logStore)
	authHandler := handler.NewAuthHandler(authService)

	s.mux.HandleFunc("/health", handler.HandleHealth)
	s.mux.HandleFunc("/ingest", ingestHandler.Handle)
	s.mux.HandleFunc("/logs", logHandler.Handle)
	s.mux.HandleFunc("/auth/register", authHandler.HandleRegister)
	s.mux.HandleFunc("/auth/login", authHandler.HandleLogin)
	s.mux.HandleFunc("/auth/logout", authHandler.HandleLogout)
}
