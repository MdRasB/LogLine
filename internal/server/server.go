package server

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/MdRasB/LogLine/internal/auth"
	"github.com/MdRasB/LogLine/internal/db"
	"github.com/MdRasB/LogLine/internal/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	addr         string
	mux          *http.ServeMux
	db           *pgxpool.Pool
	logStore     db.DBStore
	userStore    db.UserStore
	sessionStore db.SessionStore
	logger       *slog.Logger
	authService  *auth.Service

	authMiddleware     middleware.Middleware
	loggingMiddleware  middleware.Middleware
	recoveryMiddleware middleware.Middleware
	ratelimiter        *middleware.RateLimiter
}

func NewServer(addr, dbstore string, requestPerSecond float64, burst int) *Server {
	mux := http.NewServeMux()

	pool, err := db.New(dbstore)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
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
		authMiddleware:     authMiddleware,
		loggingMiddleware:  loggingMiddleware,
		recoveryMiddleware: recoveryMiddleware,
		ratelimiter:        rateLimiter,
	}

	s.registerRoutes()

	return s
}

func (s *Server) Start() error {
	log.Printf("Server running on %s\n", s.addr)
	return http.ListenAndServe(s.addr, s.mux)
}
