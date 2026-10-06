include .env
export

MIGRATE_DATABASE_URL := postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

create-migrate:
	migrate create -ext sql -dir migrations/catalog -seq init

migrate-up:
	migrate -path migrations/catalog -database "$(MIGRATE_DATABASE_URL)&x-migrations-table=catalog_migrations" up