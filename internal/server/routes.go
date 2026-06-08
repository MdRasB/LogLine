package server

import (
	"net/http"

	"github.com/MdRasB/LogLine/internal/auth"
	"github.com/MdRasB/LogLine/internal/handler"
	"github.com/MdRasB/LogLine/internal/middleware"
)

func (s *Server) registerRoutes() {

	authService := auth.NewService(
		&s.userStore,
		&s.sessionStore,
	)

	authMiddleware := middleware.AuthMiddleware(authService)

	ingestHandler := handler.NewIngestHandler(&s.logStore)
	logHandler := handler.NewLogHandler(&s.logStore)
	authHandler := handler.NewAuthHandler(authService)

	//public routes:
	s.mux.HandleFunc("/health", handler.HandleHealth)
	s.mux.HandleFunc("/auth/register", authHandler.HandleRegister)
	s.mux.HandleFunc("/auth/login", authHandler.HandleLogin)

	//protected routes:
	//s.mux.HandleFunc("/ingest", ingestHandler.Handle)
	s.mux.Handle(
		"/ingest",
		authMiddleware(
			http.HandlerFunc(ingestHandler.Handle),
		),
	)
	//s.mux.HandleFunc("/logs", logHandler.Handle)
	s.mux.Handle(
		"/logs",
		authMiddleware(
			http.HandlerFunc(logHandler.Handle),
		),
	)
	//s.mux.HandleFunc("/auth/logout", authHandler.HandleLogout)
	s.mux.Handle(
		"/auth/logout",
		authMiddleware(
			http.HandlerFunc(authHandler.HandleLogout),
		),
	)
}
