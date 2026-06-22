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
	loggingMiddleware := middleware.Logging(s.logger)

	ingestHandler := handler.NewIngestHandler(&s.logStore)
	logHandler := handler.NewLogHandler(&s.logStore)
	authHandler := handler.NewAuthHandler(authService)

	//public routes:
	//s.mux.HandleFunc("/health", handler.HandleHealth)
	s.mux.Handle("/health",
		middleware.Chain(
			http.HandlerFunc(handler.HandleHealth),
			middleware.RequestID,
			loggingMiddleware,
		),
	)

	//s.mux.HandleFunc("/auth/register", authHandler.HandleRegister)
	s.mux.Handle("/auth/register",
		middleware.Chain(
			http.HandlerFunc(authHandler.HandleRegister),
			middleware.RequestID,
			loggingMiddleware,
		),
	)

	//s.mux.HandleFunc("/auth/login", authHandler.HandleLogin)
	s.mux.Handle("/auth/login",
		middleware.Chain(
			http.HandlerFunc(authHandler.HandleLogin),
			middleware.RequestID,
			loggingMiddleware,
		),
	)

	//protected routes:
	//s.mux.HandleFunc("/ingest", ingestHandler.Handle)
	s.mux.Handle("/ingest",
		middleware.Chain(
			http.HandlerFunc(ingestHandler.Handle),
			middleware.RequestID,
			loggingMiddleware,
			authMiddleware,
		),
	)

	//s.mux.HandleFunc("/logs", logHandler.Handle)
	s.mux.Handle("/logs",
		middleware.Chain(
			http.HandlerFunc(logHandler.Handle),
			middleware.RequestID,
			loggingMiddleware,
			authMiddleware,
		),	
	)
	//s.mux.HandleFunc("/auth/logout", authHandler.HandleLogout)
	s.mux.Handle("/auth/logout",
		middleware.Chain(
			http.HandlerFunc(authHandler.HandleLogout),
			middleware.RequestID,
			loggingMiddleware,
			authMiddleware,
		),
	)
}
