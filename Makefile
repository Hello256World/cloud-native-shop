include .env
export

MIGRATE_DATABASE_URL := postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

create-migrate:
	migrate create -ext sql -dir migrations/identity -seq init

migrate-up:
	migrate -path migrations/identity -database "$(MIGRATE_DATABASE_URL)&x-migrations-table=identity_migrations" up