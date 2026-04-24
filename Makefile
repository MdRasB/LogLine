.PHONY: run build db-up db-down migrate migrate-status test db-run

run:
	go run cmd/api/main.go

build:
	go build -o logline ./cmd/api

db-up:
	docker compose up -d

db-down:
	docker compose down

migrate:
	GOOSE_DRIVER=postgres \
	GOOSE_DBSTRING="postgres://logline:password@localhost:5432/loglinedb" \
	go run github.com/pressly/goose/v3/cmd/goose \
	-dir migrations \
	up

migrate-status:
	GOOSE_DRIVER=postgres \
	GOOSE_DBSTRING="postgres://logline:password@localhost:5432/loglinedb" \
	go run github.com/pressly/goose/v3/cmd/goose \
	-dir migrations \
	status

db-run:
	docker exec -it loglinedb psql -U logline -d loglinedb

test:
	go test ./...