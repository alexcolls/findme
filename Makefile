.PHONY: help backend-dev backend-build backend-test mobile-ios mobile-android docker-up docker-down install-mobile install-backend clean

# Default target
help:
	@echo "FindMe Development Commands"
	@echo ""
	@echo "Backend:"
	@echo "  make backend-dev          - Run backend in development mode"
	@echo "  make backend-build        - Build backend binary"
	@echo "  make backend-test         - Run backend tests"
	@echo "  make install-backend      - Install backend dependencies"
	@echo ""
	@echo "Mobile:"
	@echo "  make mobile-ios           - Run iOS app"
	@echo "  make mobile-android       - Run Android app"
	@echo "  make install-mobile       - Install mobile dependencies"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-up            - Start all services (PostgreSQL, Qdrant, Redis)"
	@echo "  make docker-down          - Stop all services"
	@echo "  make docker-clean         - Remove all volumes"
	@echo ""
	@echo "Utils:"
	@echo "  make clean                - Clean build artifacts"
	@echo "  make setup                - Complete project setup"

# Backend commands
backend-dev:
	cd backend && go run cmd/api/main.go

backend-build:
	cd backend && go build -o bin/api cmd/api/main.go

backend-test:
	cd backend && go test -v ./...

backend-test-coverage:
	cd backend && go test -cover ./...

install-backend:
	cd backend && go mod download

# Mobile commands
mobile-ios:
	cd mobile && npm run ios

mobile-android:
	cd mobile && npm run android

mobile-start:
	cd mobile && npm start

install-mobile:
	cd mobile && npm install
	cd mobile/ios && pod install

# Docker commands
docker-up:
	cd backend && docker-compose up -d

docker-down:
	cd backend && docker-compose down

docker-logs:
	cd backend && docker-compose logs -f

docker-clean:
	cd backend && docker-compose down -v

docker-rebuild:
	cd backend && docker-compose up -d --build

# Setup commands
setup: install-backend install-mobile docker-up
	@echo "✅ Setup complete!"
	@echo ""
	@echo "Next steps:"
	@echo "1. Configure backend/.env"
	@echo "2. Run 'make backend-dev' to start backend"
	@echo "3. Run 'make mobile-ios' or 'make mobile-android'"

# Clean commands
clean:
	rm -rf backend/bin
	rm -rf mobile/node_modules
	rm -rf mobile/ios/Pods
	cd backend && go clean

# Git commands
commit:
	git add .
	git status

push:
	git push origin main
