# syntax=docker/dockerfile:1

# build the app
FROM golang:1.22.5 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /cats-social cmd/api/main.go

# deploy the app binary into a lean image
FROM gcr.io/distroless/static-debian12

WORKDIR /

COPY --from=build /cats-social /cats-social

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT [ "/cats-social" ]
