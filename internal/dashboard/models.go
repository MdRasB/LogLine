// Package dashboard is the module for dashboard html (webpage interface)
package dashboard

import "github.com/MdRasB/LogLine/internal/model"

type DashboardData struct {
	Logs       []model.Logs
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
