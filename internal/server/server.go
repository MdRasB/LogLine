package server

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/MdRasB/LogLine/internal/auth"
	"github.com/MdRasB/LogLine/internal/db"
	"github.com/MdRasB/LogLine/internal/middleware"
	"github.com/MdRasB/LogLine/internal/web"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	addr       string
	mux        *http.ServeMux
	httpServer *http.Server

	db           *pgxpool.Pool
	logStore     db.DBStore
	userStore    db.UserStore
	sessionStore db.SessionStore
	logger       *slog.Logger
	authService  *auth.Service
	startedAt    time.Time
	version      string

	templates *web.TemplateManager

	authMiddleware     middleware.Middleware
	loggingMiddleware  middleware.Middleware
	recoveryMiddleware middleware.Middleware
	ratelimiter        *middleware.RateLimiter
}

func NewServer(addr, dbstore string, requestPerSecond float64, burst int, startedAt time.Time, version string) (*Server, error) {
	mux := http.NewServeMux()

	pool, err := db.New(dbstore)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	dbStore := db.NewLogStore(pool)
	usrStore := db.NewUserStore(pool)
	sessnStore := db.NewSessionStore(pool)
	logger := slog.New(slog.NewTextHandler(log.Writer(), nil))

	authService := auth.NewService(
		usrStore,
		sessnStore,
	)

	rateLimiter := middleware.NewRateLimiter(
		requestPerSecond, // request per second
		burst,            // burst

	)

	templates, err := web.NewTemplateManager()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize template manager: %w", err)
	}

	// Middleware Variables
	authMiddleware := middleware.AuthMiddleware(authService)
	loggingMiddleware := middleware.Logging(logger)
	recoveryMiddleware := middleware.Recovery(logger)

	s := &Server{
		addr:               addr,
		mux:                mux,
		db:                 pool,
		logStore:           *dbStore,
		userStore:          *usrStore,
		sessionStore:       *sessnStore,
		logger:             logger,
		authService:        authService,
		startedAt:          startedAt,
		version:            version,
		templates:          templates,
		authMiddleware:     authMiddleware,
		loggingMiddleware:  loggingMiddleware,
		recoveryMiddleware: recoveryMiddleware,
		ratelimiter:        rateLimiter,
	}

	s.registerRoutes()

	s.httpServer = &http.Server{
		Addr:    s.addr,
		Handler: s.mux,
	}
	return s, nil
}

func (s *Server) Start() error {
	log.Printf("Server running on %s\n", s.addr)

	err := s.httpServer.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		return err
	}

	log.Println("HTTP server closed down normally")
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("shutting down the HTTP server gracefully...")

	err := s.httpServer.Shutdown(ctx)
	fmt.Println("Closing databse connection pool")
	if s.db != nil {
		s.db.Close()
	}

	if err != nil {
		return err
	}

	return nil
}
