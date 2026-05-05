package db

import (
	"fmt"
	"strings"

	"github.com/MdRasB/LogLine/internal/model"
)

func GetLogsQuery(lf model.LogFilter) (string, []interface{}) {
	args := []any{}
	where := []string{}

	if lf.Service != "" {
		args = append(args, lf.Service)
		where = append(where, fmt.Sprintf("service = $%d", len(args)))
	}

	if lf.Level != "" {
		args = append(args, lf.Level)
		where = append(where, fmt.Sprintf("level = $%d", len(args)))
	}

	if lf.Search != "" {
		args = append(args, "%"+lf.Search+"%")
		where = append(where, fmt.Sprintf("message ILIKE $%d", len(args)))
	}

	if lf.From != nil {
		args = append(args, *lf.From) 
		where = append(where, fmt.Sprintf("timestamp >= $%d", len(args)))
	}

	if lf.To != nil {
		args = append(args, *lf.To)
		where = append(where, fmt.Sprintf("timestamp <= $%d", len(args)))
	}
	
	qry := `SELECT id, level, message,service, timestamp, metadata, created_at FROM logs`

	if len(where) > 0 {
		qry += " Where " + strings.Join(where, " AND ")
	}

	offset := (lf.Page - 1) * lf.Limit
	args = append(args, lf.Limit, offset)

	qry += fmt.Sprintf(
		" ORDER BY timestamp DESC LIMIT $%d OFFSET $%d",
		len(args)-1, len(args),
	)

	return qry, args;
}


func CountLogsQuery(lf model.LogFilter) (string, []interface{}) {
	args := []any{}
	where := []string{}

	if lf.Service != "" {
		args = append(args, lf.Service)
		where = append(where, fmt.Sprintf("service = $%d", len(args)))
	}

	if lf.Level != "" {
		args = append(args, lf.Level)
		where = append(where, fmt.Sprintf("level = $%d", len(args)))
	}

	if lf.Search != "" {
		args = append(args, "%"+lf.Search+"%")
		where = append(where, fmt.Sprintf("message ILIKE $%d", len(args)))
	}

	if lf.From != nil {
		args = append(args, *lf.From) 
		where = append(where, fmt.Sprintf("timestamp >= $%d", len(args)))
	}

	if lf.To != nil {
		args = append(args, *lf.To)
		where = append(where, fmt.Sprintf("timestamp <= $%d", len(args)))
	}

	qry := "SELECT COUNT(*) FROM logs"

	if len(where) > 0 {
		qry += " Where " + strings.Join(where, " AND ")
	}

	return qry, args
}