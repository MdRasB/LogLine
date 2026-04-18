package db

import (
	"context"

	"github.com/MdRasB/LogLine/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBStore struct {
	db *pgxpool.Pool
}

func NewLogStore(db *pgxpool.Pool) *DBStore {
	return &DBStore{
		db: db,
	}
}

func (s *DBStore) Insert(log model.Logs) error {
	query := `
			Insert Into logs (level, message, service, timestamp, metadata)
			Values ($1, $2, $3, $4, $5)
	`

	_, err := s.db.Exec(
		context.Background(),
		query,
		log.Level,
		log.Message,
		log.Service,
		log.Timestamp,
		log.Metadata,
	)

	if err != nil {
		return nil
	}
	
	return nil
}
