# Simple Makefile for a Go project

# Build the application
all: build

build:
	@echo "Building..."
	@weaver generate ./...
	@go build -o ./cmd/api/gahoot cmd/api/main.go

# Run the application
run:
	@weaver generate ./...
	@SERVICEWEAVER_CONFIG=weaver.toml go run ./cmd/api/main.go

# Create DB container
docker-run:
	@if docker compose up 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up; \
	fi

# Shutdown DB container
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

# Test the application
test:
	@echo "Testing..."
	@go test ./tests -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
	    air; \
	    echo "Watching...";\
	else \
	    read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
	    if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
	        go install github.com/cosmtrek/air@latest; \
	        air; \
	        echo "Watching...";\
	    else \
	        echo "You chose not to install air. Exiting..."; \
	        exit 1; \
	    fi; \
	fi

# Dashboard
dashboard:
	@echo "Opening Dashboard..."
	@weaver single dashboard

# Migrations
migrate:
	@echo "Running migrations..."
	@go run ./cmd/migrate/main.go

.PHONY: all build run test clean watch dashboard migrate
