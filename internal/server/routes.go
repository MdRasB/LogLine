// Package server handle this project's servers and routes
package server

import (
	"net/http"

	"github.com/MdRasB/LogLine/internal/handler"
	"github.com/MdRasB/LogLine/internal/middleware"
)

func (s *Server) registerRoutes() {
	ingestHandler := handler.NewIngestHandler(&s.logStore)
	logHandler := handler.NewLogHandler(&s.logStore)
	authHandler := handler.NewAuthHandler(s.authService)

	// public routes:
	s.mux.Handle(
		"/health",
		middleware.Chain(
			http.HandlerFunc(handler.HandleHealth),
			middleware.RequestID,
			s.recoveryMiddleware,
			s.loggingMiddleware,
			s.ratelimiter.Middleware,
		),
	)

	s.mux.Handle(
		"/auth/register",
		middleware.Chain(
			http.HandlerFunc(authHandler.HandleRegister),
			middleware.RequestID,
			s.recoveryMiddleware,
			s.loggingMiddleware,
			s.ratelimiter.Middleware,
		),
	)

	s.mux.Handle(
		"/auth/login",
		middleware.Chain(
			http.HandlerFunc(authHandler.HandleLogin),
			middleware.RequestID,
			s.recoveryMiddleware,
			s.loggingMiddleware,
			s.ratelimiter.Middleware,
		),
	)

	// protected routes:
	s.mux.Handle(
		"/ingest",
		middleware.Chain(
			http.HandlerFunc(ingestHandler.Handle),
			middleware.RequestID,
			s.recoveryMiddleware,
			s.loggingMiddleware,
			s.ratelimiter.Middleware,
			s.authMiddleware,
		),
	)

	s.mux.Handle(
		"/logs",
		middleware.Chain(
			http.HandlerFunc(logHandler.Handle),
			middleware.RequestID,
			s.recoveryMiddleware,
			s.loggingMiddleware,
			s.ratelimiter.Middleware,
			s.authMiddleware,
		),
	)

	s.mux.Handle(
		"/auth/logout",
		middleware.Chain(
			http.HandlerFunc(authHandler.HandleLogout),
			middleware.RequestID,
			s.recoveryMiddleware,
			s.loggingMiddleware,
			s.ratelimiter.Middleware,
			s.authMiddleware,
		),
	)
}
