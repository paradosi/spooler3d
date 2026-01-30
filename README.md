# Spooler3D

REST API for managing 3D printing filament spools. Built with Go and PostgreSQL.

Tracks manufacturers, filament types, spools (with NFC-based UUID identifiers), and weight history from an ESP32 scale.

## Stack

- **Go 1.23** with [Gin](https://github.com/gin-gonic/gin) HTTP framework
- **PostgreSQL 16** with [sqlx](https://github.com/jmoiron/sqlx) for queries
- **Docker Compose** for local development and deployment
- **golang-migrate** for database migrations

## Quick Start

```bash
# Clone
git clone https://github.com/paradosi/spooler3d.git
cd spooler3d

# Configure
cp .env.example .env

# Start PostgreSQL
docker compose up -d postgres

# Run migrations
migrate -path migrations -database "postgres://filament:changeme@localhost:5432/spooler3d?sslmode=disable" up

# Run the API
go run ./cmd/server
```

Or run everything with Docker:

```bash
docker compose up -d
```

The API will be available at `http://localhost:8081`.

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/health` | Health check with DB ping |
| GET | `/api/manufacturers` | List manufacturers |
| POST | `/api/manufacturers` | Create manufacturer |
| GET | `/api/manufacturers/:id` | Get manufacturer |
| PUT | `/api/manufacturers/:id` | Update manufacturer |
| DELETE | `/api/manufacturers/:id` | Delete manufacturer |
| GET | `/api/filament-types` | List filament types |
| POST | `/api/filament-types` | Create filament type |
| GET | `/api/filament-types/:id` | Get filament type |
| PUT | `/api/filament-types/:id` | Update filament type |
| DELETE | `/api/filament-types/:id` | Delete filament type |
| GET | `/api/spools` | List spools (filterable) |
| POST | `/api/spools` | Create spool |
| GET | `/api/spools/:id` | Get spool |
| PUT | `/api/spools/:id` | Update spool |
| DELETE | `/api/spools/:id` | Delete spool |
| POST | `/api/spools/:uid/weight` | Update weight by UUID (ESP32) |
| GET | `/api/spools/:id/weight-history` | Weight history for a spool |
| GET | `/api/stats` | Dashboard aggregates |

### Spool Filtering

The `GET /api/spools` endpoint supports query parameters:

- `manufacturer_id` — filter by manufacturer
- `filament_type_id` — filter by filament type
- `location` — partial match (case-insensitive)

### ESP32 Weight Update

The ESP32 reads a UUID from an NFC tag on the spool and posts the measured weight:

```bash
curl -X POST http://localhost:8081/api/spools/<uuid>/weight \
  -H "Content-Type: application/json" \
  -d '{"weight": 847.5}'
```

This atomically updates the spool's `current_weight` and inserts a record into `weight_history`.

## Database Schema

Four tables managed via migrations:

- **manufacturers** — name, website
- **filament_types** — name, print/bed temperature ranges (seeded with PLA, PETG, ABS, TPU, ASA)
- **spools** — UUID, color, diameter, weights, location, purchase info, linked to manufacturer and filament type
- **weight_history** — timestamped weight readings per spool

`remaining_weight` is computed at query time as `current_weight - spool_weight`.

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | `filament` | Database user |
| `DB_PASSWORD` | `changeme` | Database password |
| `DB_NAME` | `spooler3d` | Database name |
| `DB_SSLMODE` | `disable` | SSL mode |
| `SERVER_PORT` | `8081` | API listen port |
| `GIN_MODE` | `debug` | Gin mode (`debug` or `release`) |

## License

MIT
