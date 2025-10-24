# Makefile

# Load env vars from .env file and run the app
run:
	@echo "Running app with .env variables loaded..."
	@set -a && . .env && set +a && go run cmd/main.go

# Build binary with env loaded (if you want to test build)
build:
	@echo "Building app with .env variables loaded..."
	@set -a && . .env && set +a && go build -o app

# Clean built binary
clean:
	rm -f app

# Show env vars loaded from .env for debug
env:
	@set -a && . .env && set +a && env | grep -E 'DB_|JWT_|ACCESS_TOKEN_TTL|REFRESH_TOKEN_TTL'

.PHONY: run build clean env
