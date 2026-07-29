package db

import (
	"context"

	"github.com/MdRasB/LogLine/internal/dashboard"
)

func (s *DBStore) GetDashboardStats(ctx context.Context) (*dashboard.DashboardStats, error) {
	stats := &dashboard.DashboardStats{}

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

	rows, err := s.db.Query(
		ctx,
		`
		SELECT TO_CHAR(DATE_TRUNC('hour', timestamp),'YYYY-MM-DD HH24:00') AS hour, COUNT(*)
		FROM logs
		GROUP BY DATE_TRUNC('hour', timestamp)
		ORDER BY DATE_TRUNC('hour', timestamp);
		`,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var h dashboard.HourlyStat

		err := rows.Scan(&h.Hour, &h.Count)
		if err != nil {
			return nil, err
		}

		stats.Volume = append(stats.Volume, h)
	}

	rows, err = s.db.Query(
		ctx,
		`
		SELECT level, COUNT(*) FROM logs
		GROUP BY level
		ORDER BY level;
		`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var l dashboard.LevelStat

		err := rows.Scan(&l.Level, &l.Count)
		if err != nil {
			return nil, err
		}

		stats.Levels = append(stats.Levels, l)
	}

	rows, err = s.db.Query(
		ctx,
		`
		SELECT service, COUNT(*) FROM logs
		GROUP BY service
		ORDER BY service;
		`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s dashboard.ServiceStat

		err := rows.Scan(&s.Service, &s.Count)
		if err != nil {
			return nil, err
		}

		stats.Services = append(stats.Services, s)
	}

	return stats, nil
}
