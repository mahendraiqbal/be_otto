.PHONY: docker-up docker-down migrate-up migrate-down

# Docker commands
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

# Get GOPATH
GOPATH=$(shell go env GOPATH)

# Database migrations
migrate-up:
	$(GOPATH)/bin/goose -dir migrations postgres "postgres://postgres:postgres@localhost:5433/voucher_db?sslmode=disable" up

migrate-down:
	$(GOPATH)/bin/goose -dir migrations postgres "postgres://postgres:postgres@localhost:5433/voucher_db?sslmode=disable" down

# Run application
run:
	go run cmd/api/main.go

# Install dependencies
install:
	go mod download
	go install github.com/pressly/goose/v3/cmd/goose@latest
