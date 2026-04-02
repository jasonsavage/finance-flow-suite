# Finance Flow API

Backend REST API for the Finance Flow application written in Go.

## Getting Started

Make sure docker is running, then run the following command to start the database and API:

```bash
    docker compose up
    # or
    docker compose up -d
```

To run the API without docker:

```bash
docker compose up db

    go run cmd/web/main.go
```

## Testing the API

```bash
    ./scripts/register.sh to register a new user
    ./scripts/login.sh to login and get a token
```
