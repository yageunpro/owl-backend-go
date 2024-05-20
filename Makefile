.PHONY: build-dev build-prod vendor

build: build-dev

build-dev:
	@docker build . -t owl-backend-go -f ./cmd/Dockerfile --target dev

build-prod:
	@docker build . -t owl-backend-go -f ./cmd/Dockerfile --target prod

vendor:
	@go mod tidy
	@go mod vendor

env:
	@echo 'export GOOSE_DRIVER=postgres' >> .env
	@echo 'export GOOSE_DBSTRING="host=host password=password user=user dbname=name search_path=public sslmode=disable"' >> .env
	@echo 'export GOOSE_MIGRATION_DIR="./asset/sql/schema"' >> .env

migration:
	@echo '#!/bin/bash\n\nif [ -f .env ]; then\n	source .env\nfi\n\ngoose -s -table migration_history $$@\n' > migrate.sh
	@chmod +x migrate.sh