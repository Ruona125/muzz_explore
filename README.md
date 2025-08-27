# Explore_Muzz Service (Go + gRPC + Postgres)

A small, fast gRPC service written in Go with PostgreSQL for storage. Ships with Docker and docker-compose for a zero-friction spin-up.

## Tech Stack

* **Language:** Go
* **RPC:** gRPC (Protocol Buffers)
* **Database:** PostgreSQL
* **Containerization:** Docker, docker-compose
* **Config:** `.env`

## Quick Start

### 1) Run with Docker Compose (recommended)

```bash
# From the repo root
docker-compose up --build
# (Compose V2 users can run: docker compose up --build)
```

This will:

* Build the Go service image
* Start Postgres
* Start the gRPC service and run DB migrations on startup

### 2) Run Locally (no Docker)

Prereqs: Go installed, local PostgreSQL running and reachable via your `.env`.

```bash
# 1) Create the database so migrations can run
createdb explore           # or create it via your DB client of choice

# 2) Install dependencies
go mod tidy

# 3) Run the service
go run ./cmd/explore-service
```

### 3) Run Tests

```bash
go test -run TestExplorerServer -v ./test
```

## Environment Variables

Configuration is read from `.env`. Open it to see exact names and defaults used by the service. Typical variables include:

* `POSTGRES_HOST`, `POSTGRES_PORT`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`
* Or a single `DATABASE_URL` connection string
* `GRPC_PORT` (the port the service listens on)

> **Note:** Defaults live in `.env`. If you change ports or creds, update both your local env and `docker-compose.yml` to match.

## gRPC Usage

* **Protos:** See the `proto/` directory (or where your `.proto` files live in this repo).
* **Sample call with grpcurl** (if reflection is enabled in the server):

  ```bash
  grpcurl -plaintext localhost:${GRPC_PORT:-50051} list
  # Replace with an actual method:
  grpcurl -plaintext -d '{}' localhost:${GRPC_PORT:-50051} package.Service/Method
  ```

If reflection isn’t enabled, generate a client from `proto/` and call methods directly.

## Database & Migrations

* The app **expects a database named `explore`** when running locally, so the migration logic can create tables automatically on startup.
* With Docker Compose, the DB is created for you (check `docker-compose.yml`).
* If you need to reset locally: drop and recreate the `explore` DB, then restart the service.

## Project Structure (high-level)

```
.
├─ cmd/
│  └─ explore-service/      # main entrypoint (go run ./cmd/explore-service)
├─ internal/                # app/internal packages (business logic, adapters, repos)
├─ proto/                   # .proto definitions
├─ test/                    # tests (go test -run TestExplorerServer -v ./test)
├─ docker/                  # Dockerfile(s) or related assets (if present)
├─ docker-compose.yml
├─ .env                     # local config (checked in intentionally; see note below)
└─ README.md
```

*(If your layout differs, adjust the paths above—commands in this README use what you provided.)*

## Troubleshooting

* **Port already in use:** Change `GRPC_PORT` in `.env` and/or the compose file.
* **DB connection errors:** Confirm Postgres is up, credentials match `.env`, and the `explore` database exists (for local runs).
* **Migrations didn’t run:** Confirm the service logs on startup; ensure the DB is reachable and the `explore` database exists before starting the app.
* **Compose not found:** On newer Docker installs the command is `docker compose` (space), not `docker-compose` (dash).

## Security & Why the `.env` Is Committed

Yes, the `.env` file is in the repo and the database credentials are visible. That is **intentional** for this project, for two reasons:

1. **Localhost only:** The service and database are designed to run on your machine/inside local Docker. The credentials aren’t exposed on the public internet.
2. **Zero-friction onboarding:** Anyone should be able to clone the repo and run the app without fighting secret management on day one.

> **Hard truth:** This is fine for local/dev. It’s **not** fine for staging/production. Never ship real secrets in git. Use a proper secrets manager (Vault, SSM, Doppler, 1Password, Kubernetes secrets, etc.). If these creds ever escape their intended local context, rotate them immediately.

## Make It Yours

* Want to change ports? Update `.env` and `docker-compose.yml`.
* Want to run in a container only? Stick to compose.
* Want to plug in a different Postgres? Point `DATABASE_URL` (or the discrete vars) at it.

## Contributing

* Fork, branch, commit with clear messages, open a PR.
* Keep changes small and tested (`go test -run TestExplorerServer -v ./test`).

If you want, I can tailor this to your exact `.proto` package/method names and the actual env var keys you’re using—just paste those in and I’ll wire them into the README.
