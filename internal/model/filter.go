package model

import "time"

type LogFilter struct {
	Service string
	Level   string
	Search  string
	From    *time.Time
	To      *time.Time
	Page    int
	Limit   int
}

type LogEntry struct {
	ID string
	Level     string
	Message   string
	Service   string
	Timestamp time.Time
	Metadata  map[string]interface{}
	CreatedAt time.Time
}

type PaginatedLogs struct {
	Logs	   []LogEntry
	Total   int
	Page    int
	Limit   int	
	HasMore  bool
}
