.PHONY: all
all: migrate-up run

.PHONY: build
build:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/cats-social cmd/api/main.go

.PHONY: run
run: build
	@./bin/cats-social

.PHONY: create-migration
create-migration:
	@migrate create -ext sql -dir db/migrations $(MIGRATE_NAME)

.PHONY: migrate-up
migrate-up:
	@migrate -database "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?$(DB_PARAMS)" -path db/migrations -verbose up

.PHONY: migrate-down
migrate-down:
	@migrate -database "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?$(DB_PARAMS)" -path db/migrations -verbose down

.PHONY: migrate-drop
migrate-drop:
	@migrate -database "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?$(DB_PARAMS)" -path db/migrations -verbose drop

.PHONY: migrate-version
migrate-version:
	@migrate -database "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?$(DB_PARAMS)" -path db/migrations version

.PHONY: migrate-force
migrate-force:
	@migrate -database "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?$(DB_PARAMS)" -path db/migrations force $(MIGRATE_VERSION)

.PHONY: migrate-clean
clean: migrate-down
	@rm -rf bin/

.PHONY: docker-up
docker-up:
	@docker compose up --build -d

.PHONY: docker-down
docker-down:
	@docker compose stop && docker compose down

.PHONY: docker-down-volumes
docker-down-volumes:
	@docker compose stop && docker compose down --volumes
