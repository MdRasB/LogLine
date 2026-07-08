// Package server handle this project's servers and routes
package server

import (
	"net/http"

	"github.com/MdRasB/LogLine/internal/handler"
	"github.com/MdRasB/LogLine/internal/middleware"
)

func (s *Server) registerPublicRoutes() {
	authHandler := handler.NewAuthHandler(s.authService)

	s.mux.Handle(
		"/health",
		s.publicChain(
			http.HandlerFunc(handler.HandleHealth),
		),
	)

	s.mux.Handle(
		"/auth/register",
		s.publicChain(
			http.HandlerFunc(authHandler.HandleRegister),
		),
	)

	s.mux.Handle(
		"/auth/login",
		s.publicChain(
			http.HandlerFunc(authHandler.HandleLogin),
		),
	)
}

func (s *Server) registerProtectedRoutes() {
	ingestHandler := handler.NewIngestHandler(&s.logStore)
	logHandler := handler.NewLogHandler(&s.logStore)
	authHandler := handler.NewAuthHandler(s.authService)

	s.mux.Handle(
		"/ingest",
		s.protectedChain(
			http.HandlerFunc(ingestHandler.Handle),
		),
	)

	s.mux.Handle(
		"/logs",
		s.protectedChain(
			http.HandlerFunc(logHandler.Handle),
		),
	)

	s.mux.Handle(
		"/auth/logout",
		s.protectedChain(
			http.HandlerFunc(authHandler.HandleLogout),
		),
	)
}

func (s *Server) publicChain(h http.Handler) http.Handler {
	return middleware.Chain(
		h,
		middleware.RequestID,
		s.recoveryMiddleware,
		s.loggingMiddleware,
		s.ratelimiter.Middleware,
	)
}

func (s *Server) protectedChain(h http.Handler) http.Handler {
	return middleware.Chain(
		h,
		middleware.RequestID,
		s.recoveryMiddleware,
		s.loggingMiddleware,
		s.ratelimiter.Middleware,
		s.authMiddleware,
	)
}

func (s *Server) registerRoutes() {
	s.registerPublicRoutes()
	s.registerProtectedRoutes()
}
