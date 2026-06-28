// Package server handle this project's servers and routes
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

	rateLimiter := middleware.NewRateLimiter(
		5,  // request per second
		10, // burst

	)

	// Middleware Variables
	authMiddleware := middleware.AuthMiddleware(authService)
	loggingMiddleware := middleware.Logging(s.logger)
	recoveryMiddleware := middleware.Recovery(s.logger)

	ingestHandler := handler.NewIngestHandler(&s.logStore)
	logHandler := handler.NewLogHandler(&s.logStore)
	authHandler := handler.NewAuthHandler(authService)

	// public routes:
	s.mux.Handle(
		"/health",
		middleware.Chain(
			http.HandlerFunc(handler.HandleHealth),
			middleware.RequestID,
			recoveryMiddleware,
			loggingMiddleware,
			rateLimiter.Middleware,
		),
	)

	s.mux.Handle(
		"/auth/register",
		middleware.Chain(
			http.HandlerFunc(authHandler.HandleRegister),
			middleware.RequestID,
			recoveryMiddleware,
			loggingMiddleware,
			rateLimiter.Middleware,
		),
	)

	s.mux.Handle(
		"/auth/login",
		middleware.Chain(
			http.HandlerFunc(authHandler.HandleLogin),
			middleware.RequestID,
			recoveryMiddleware,
			loggingMiddleware,
			rateLimiter.Middleware,
		),
	)

	// protected routes:
	s.mux.Handle(
		"/ingest",
		middleware.Chain(
			http.HandlerFunc(ingestHandler.Handle),
			middleware.RequestID,
			recoveryMiddleware,
			loggingMiddleware,
			rateLimiter.Middleware,
			authMiddleware,
		),
	)

	s.mux.Handle(
		"/logs",
		middleware.Chain(
			http.HandlerFunc(logHandler.Handle),
			middleware.RequestID,
			recoveryMiddleware,
			loggingMiddleware,
			rateLimiter.Middleware,
			authMiddleware,
		),
	)

	s.mux.Handle(
		"/auth/logout",
		middleware.Chain(
			http.HandlerFunc(authHandler.HandleLogout),
			middleware.RequestID,
			recoveryMiddleware,
			loggingMiddleware,
			rateLimiter.Middleware,
			authMiddleware,
		),
	)
}
