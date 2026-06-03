package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	connString :=
		"postgres://postgres:postgres@localhost:5432/logline_test?sslmode=disable"

	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		t.Fatalf("connect test db: %v", err)
	}

	return pool
}