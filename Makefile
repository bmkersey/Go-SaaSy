# Makefile for SaaSy
include .env
# Run the app and database with build
up:
	docker-compose up --build

# Start app and database without rebuilding
start:
	docker-compose up

# Shut everything down
down:
	docker-compose down

# Run Goose migrations in isolated container
migrate:
	docker-compose --profile migrate -f docker-compose.yml -f docker-compose.migrate.yml  run --rm migrate

# Run goose down migrations in isolated container
migrate-down:
	docker-compose -f docker-compose.yml -f docker-compose.migrate.yml run --rm migrate goose -dir sql/schema postgres "$(DB_URL)" down

# Build the API binary manually (local dev)
build:
	go build -o saasy ./cmd/saasy

# Run unit tests
test:
	go test ./...

# Format Go code
fmt:
	go fmt ./...

# View running containers
ps:
	docker-compose ps

# Rebuild everything cleanly
rebuild:
	docker-compose down
	docker-compose build
	docker-compose up
# Sends the test email
reset-email:
	docker compose run --rm app go run cmd/sendreset/main.go
# Tail logs
logs:
	docker-compose logs -f --tail=100

.PHONY: up start down migrate build test fmt ps rebuild logs
