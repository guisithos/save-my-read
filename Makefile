.PHONY: run test clean db-init db-migrate db-rollback

# Load environment variables from .env file
include .env
export

run:
	@echo "Starting the application..."
	go run cmd/api/main.go

test:
	@echo "Running tests..."
	@go test ./...

clean:
	@echo "Cleaning..."
	@rm -f tmp/* 

dev:
	@echo "Starting development server..."
	go run cmd/api/main.go 

db-init:
	@echo "Initializing database..."
	@go run scripts/db/init.go

db-migrate:
	@echo "Running migrations..."
	@migrate -path migrations -database "$(DATABASE_URL)" up

db-rollback:
	@echo "Rolling back last migration..."
	@migrate -path migrations -database "$(DATABASE_URL)" down 1 