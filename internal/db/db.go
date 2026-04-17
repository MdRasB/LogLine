package db

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(dbStr string) (*pgxpool.Pool, error)