// Package db controlls the database logic for this logline project
package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New(dbStr string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), dbStr)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
