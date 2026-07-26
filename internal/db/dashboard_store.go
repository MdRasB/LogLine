package db

import (
	"context"

	"github.com/MdRasB/LogLine/internal/dashboard"
)

func (s *DBStore) GetDashboardStats(
	ctx context.Context,
) (*dashboard.DashboardStats, error) {
	stats := &dashboard.DashboardStats{}

	// ----------------------------------------
	// Total Logs
	// ----------------------------------------

	err := s.db.QueryRow(
		ctx,
		`
		SELECT COUNT(*)
		FROM logs
		`,
	).Scan(&stats.TotalLogs)
	if err != nil {
		return nil, err
	}

	// ----------------------------------------
	// Total Services
	// ----------------------------------------

	err = s.db.QueryRow(
		ctx,
		`
		SELECT COUNT(DISTINCT service)
		FROM logs
		`,
	).Scan(&stats.TotalServices)
	if err != nil {
		return nil, err
	}

	// ----------------------------------------
	// Error Logs
	// ----------------------------------------

	err = s.db.QueryRow(
		ctx,
		`
		SELECT COUNT(*)
		FROM logs
		WHERE level = 'error'
		`,
	).Scan(&stats.ErrorLogs)
	if err != nil {
		return nil, err
	}

	return stats, nil
}
