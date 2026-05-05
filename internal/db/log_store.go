package db

import (
	"context"
	"encoding/json"
	"fmt"

	//"fmt"

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
		//fmt.Errorf("db: failed to insert log: %v", err)
		return err
	}
	
	return nil
}

func (s *DBStore) GetLogs(
	lf model.LogFilter,
	ctx context.Context,
)(*model.PaginatedLogs, error) {
	countQ, countArgs := CountLogsQuery(lf)

	var total int
	err := s.db.QueryRow(ctx, countQ, countArgs...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("db: count logs: %w", err)
	}

	query, args := GetLogsQuery(lf)

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("db: query logs: %w", err)
	}

	defer rows.Close()

	logs := []model.LogEntry{}

	for rows.Next() {
		var l model.LogEntry
		var metaBytes []byte

		err := rows.Scan(
			&l.ID,
			&l.Level,
			&l.Message,
			&l.Service,
			&l.Timestamp,
			&metaBytes,
			&l.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("db: scan log row: %w", err)
		}

		if metaBytes != nil {
			_ = json.Unmarshal(metaBytes, &l.Metadata)
		}

		logs = append(logs, l)
		
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("db: rows error: %w", err)
	}

	return &model.PaginatedLogs{
		Logs: logs,
		Total: total,
		Page: lf.Page,
		Limit: lf.Limit,
		HasMore: (lf.Page * lf.Limit) < total,
	}, nil
}
