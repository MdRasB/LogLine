.PHONY: \
	run build compile \
	db-up db-down db-logs db-run \
	migrate migrate-status \
	test check \
	docker-build docker-up docker-ps \
	docker-shell docker-rebuild docker-clean \
	compose-build compose-up compose-down compose-logs

# Go application

run:
	go run cmd/api/main.go

build:
	go build -o logline ./cmd/api

compile:
	go build ./...

# Database

db-up:
	docker compose up -d db

db-down:
	docker compose stop db

db-logs:
	docker compose logs -f db

db-run:
	docker exec -it loglinedb psql -U logline -d loglinedb

# Database migrations

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

# Testing / linting

test:
	go test ./...

check:
	golangci-lint run -v

# Docker

docker-build:
	docker build -t logline:dev .

docker-up:
	docker compose up -d --build

docker-ps:
	docker compose ps

docker-shell:
	docker compose exec app /bin/sh

docker-rebuild:
	docker compose build --no-cache

docker-clean:
	docker compose down --rmi local


#  Docker compose

compose-build:
	docker compose build

compose-up:
	docker compose up -d

compose-down:
	docker compose down

compose-logs:
	docker compose logs -f
