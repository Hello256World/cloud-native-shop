include .env
export

MIGRATE_DATABASE_URL := postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

# Order matters: platform FIRST (citext extension + platform.set_updated_at()).
MODULES     := platform identity catalog inventory cart orders payment
MODULES_REV := payment orders cart inventory catalog identity platform

# Only modules whose migrations folder exists (payment has none yet).
ACTIVE     := $(foreach m,$(MODULES),$(if $(wildcard migrations/$(m)),$(m)))
ACTIVE_REV := $(foreach m,$(MODULES_REV),$(if $(wildcard migrations/$(m)),$(m)))

# Pinned in go.mod via `go get -tool`; postgres driver needs the build tag.
MIGRATE := go run -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate

migrate_cmd = $(MIGRATE) -path migrations/$(1) \
	-database "$(MIGRATE_DATABASE_URL)&x-migrations-table=$(1)_migrations"


.PHONY: migrate-up migrate-down-all migrate-up-module migrate-down-module \
        migrate-version migrate-force migrate-new

## Apply all modules, in dependency order
migrate-up:
	@set -e; for m in $(ACTIVE); do \
		echo "==> up: $$m"; \
		$(call migrate_cmd,$${m}) up; \
	done

## Roll back EVERYTHING, in reverse order (dev only!)
migrate-down-all:
	@set -e; for m in $(ACTIVE_REV); do \
		echo "==> down: $$m"; \
		$(call migrate_cmd,$${m}) down -all; \
	done

## make migrate-up-module MODULE=orders
migrate-up-module:
	@test -n "$(MODULE)" || (echo "usage: make $@ MODULE=<name>"; exit 1)
	$(call migrate_cmd,$(MODULE)) up

## make migrate-down-module MODULE=orders   (rolls back ONE step)
migrate-down-module:
	@test -n "$(MODULE)" || (echo "usage: make $@ MODULE=<name>"; exit 1)
	$(call migrate_cmd,$(MODULE)) down 1

## Show current version of every module
migrate-version:
	@for m in $(ACTIVE); do \
		printf "%-10s " $$m; \
		$(call migrate_cmd,$${m}) version 2>&1 | tail -1; \
	done

## make migrate-force MODULE=orders VERSION=1   (fix a "dirty" state)
migrate-force:
	@test -n "$(MODULE)" -a -n "$(VERSION)" || (echo "usage: make $@ MODULE=<name> VERSION=<n>"; exit 1)
	$(call migrate_cmd,$(MODULE)) force $(VERSION)

## make migrate-new MODULE=catalog NAME=add_brand
migrate-new:
	@test -n "$(MODULE)" -a -n "$(NAME)" || (echo "usage: make $@ MODULE=<name> NAME=<name>"; exit 1)
	$(MIGRATE) create -ext sql -dir migrations/$(MODULE) -seq $(NAME)

create-migrate:
	migrate create -ext sql -dir migrations/payment -seq init

get-migrate:
	go get -tool github.com/golang-migrate/migrate/v4/cmd/migrate@latest


.PHONY: sqlc
sqlc:
	go tool sqlc generate

.PHONY: swag
## Generate OpenAPI docs into ./docs
swag:
	go tool swag init -g cmd/api/main.go --parseInternal


.PHONY: lint
lint:
	golangci-lint run ./...


random-base64:
	openssl rand -base64 48