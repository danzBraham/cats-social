build:
	@go build -o bin/cats-social cmd/server/main.go

run: build
	@./bin/cats-social