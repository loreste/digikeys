.PHONY: dev-up dev-down backend-run backend-test web-dev mobile-dev clean

# Backend
backend-build:
	cd backend && go build -o bin/server cmd/server/main.go

backend-run:
	cd backend && go run cmd/server/main.go serve

backend-test:
	cd backend && go test ./... -cover

backend-migrate:
	cd backend && go run cmd/server/main.go migrate

backend-seed:
	cd backend && go run cmd/server/main.go seed

# Frontend
web-dev:
	cd web && npm run dev

web-build:
	cd web && npm run build

# Mobile
mobile-dev-android:
	cd mobile && npx react-native run-android

mobile-dev-ios:
	cd mobile && npx react-native run-ios

# Development
dev-up:
	docker compose up -d postgres minio
	@echo "Waiting for PostgreSQL..."
	@sleep 3
	@make backend-migrate
	@echo "DIGIKEYS dev stack ready"

dev-down:
	docker compose down

# Tests
test: backend-test

# Clean
clean:
	rm -rf backend/bin web/.next web/node_modules mobile/node_modules
