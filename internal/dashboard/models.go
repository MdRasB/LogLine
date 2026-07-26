// Package dashboard contains models used by the HTML dashboard.
package dashboard

import "github.com/MdRasB/LogLine/internal/model"

type DashboardData struct {
	Logs       []model.LogEntry
	Stats      DashboardStats
	Filters    DashboardFilters
	Pagination PaginationData
}

type DashboardFilters struct {
	Service string
	Level   string
	Search  string
}

type PaginationData struct {
	CurrentPage int
	TotalPages  int
	HasNext     bool
	HasPrev     bool
}

type DashboardStats struct {
	TotalLogs     int
	TotalServices int
	ErrorLogs     int

	Volume   []HourlyStat
	Levels   []LevelStat
	Services []ServiceStat
}

type HourlyStat struct {
	Hour  string
	Count int
}

type LevelStat struct {
	Level string
	Count int
}

type ServiceStat struct {
	Service string
	Count   int
}
