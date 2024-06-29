all: migrate-up run

build:
	@go build -o bin/cats-social cmd/api/main.go

run: build
	@./bin/cats-social

create-migration:
	@migrate create -ext sql -dir db/migrations $(MIGRATE_NAME)

migrate-up:
	@migrate -database "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?$(DB_PARAMS)" -path db/migrations -verbose up

migrate-down:
	@migrate -database "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?$(DB_PARAMS)" -path db/migrations -verbose down

migrate-drop:
	@migrate -database "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?$(DB_PARAMS)" -path db/migrations -verbose drop

migrate-version:
	@migrate -database "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?$(DB_PARAMS)" -path db/migrations version

migrate-force:
	@migrate -database "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?$(DB_PARAMS)" -path db/migrations force $(MIGRATE_VERSION)

clean: migrate-down
	@rm -rf bin/

.PHONY:
	all build run create-migration migrate-up migrate-down migrate-drop clean
