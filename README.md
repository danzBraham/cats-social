# Cats Social

An application where cat owners can match their cats with each other.

## Tech Stack

[![Go Badge](https://img.shields.io/badge/golang-golang?style=for-the-badge&logo=go&logoColor=go&logoSize=auto&labelColor=auto&color=black)](https://go.dev/)
[![PostgreSQL Badge](https://img.shields.io/badge/postgresql-postgresql?style=for-the-badge&logo=postgresql&logoColor=postgresql&logoSize=auto&labelColor=auto&color=black)](https://www.postgresql.org/)
[![Docker Badge](https://img.shields.io/badge/docker-docker?style=for-the-badge&logo=docker&logoColor=docker&logoSize=auto&labelColor=auto&color=black)](https://www.docker.com/)

## Features

- **Authentication & Authorization**: User registration and login.
- **Managing Cats**: CRUD operations for cats.
- **Matching Cats**: Matchmaking and managing cat matches.

## Run Locally

Prerequisites

- Docker

Clone the project

```bash
git clone https://github.com/danzBraham/cats-social
```

Go to the project directory

```bash
cd cats-social
```

Create an .env file and add the following environment variables to your .env file

```env
export DB_USERNAME=
export DB_PASSWORD=
export DB_HOST=
export DB_PORT=
export DB_NAME=
export DB_PARAMS=sslmode=disable
export JWT_SECRET=
export BCRYPT_SALT=8 # don't use 8 in prod! use > 10
```

Run docker

```bash
docker compose up --build -d
```

## API Reference

For a detailed API reference, check out the [API Reference](./api-reference.md).

## Feedback

If you have any feedback, please dm me at
[![x](https://img.shields.io/badge/danzBraham-x?style=for-the-badge&logo=x&logoColor=x&logoSize=auto&labelColor=auto&color=black)](https://x.com/danzBraham)
