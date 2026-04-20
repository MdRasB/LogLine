package server

import (
	"log"
	"net/http"

	"github.com/MdRasB/LogLine/internal/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	addr string
	mux  *http.ServeMux
	db   *pgxpool.Pool
	logStore db.DBStore	
}

func NewServer(addr, dbstore string) *Server {
	mux := http.NewServeMux()

	pool, err := db.New(dbstore)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	dbStore := db.NewLogStore(pool)

	s := &Server{
		addr: addr,
		mux:  mux,
		db:   pool,
		logStore: *dbStore,
	}

	s.registerRoutes()

	return s
}

func (s *Server) Start() error {
	log.Printf("Server running on %s\n", s.addr)
	return http.ListenAndServe(s.addr, s.mux)
}
