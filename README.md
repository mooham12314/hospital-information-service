# hospital-information-service

Backend service for Hospital Middleware using Go, Gin, PostgreSQL, Docker, and Nginx.

## Prerequisites

- Go 1.25+
- Docker + Docker Compose

## Quick Start (Docker)

1. Start services:

```bash
docker compose up --build -d
```

2. Check health:

```bash
curl http://localhost:8088/health
```

Expected response:

```json
{"db":"up","status":"ok"}
```

3. Stop services:

```bash
docker compose down
```

## Local Run (without Docker)

1. Ensure PostgreSQL is running and update environment variables in `.env`.
2. Run server:

```bash
go run ./cmd/server
```

3. Health check:

```bash
curl http://localhost:8080/health
```

## Project Structure

```text
cmd/server/                # application entrypoint
internal/config/           # environment config
internal/database/         # postgres connection pool
internal/http/             # gin router
internal/http/handler/     # HTTP handlers
db/init/                   # postgres initialization SQL
nginx/                     # nginx reverse proxy config
```

## Available Endpoints

- `GET /health`
- `GET /api/v1/health`
- `POST /staff/create`
- `POST /staff/login`
- `POST /patient/search` (requires `Authorization: Bearer <token>`)
- `POST /api/v1/staff/create`
- `POST /api/v1/staff/login`
- `POST /api/v1/patient/search` (requires `Authorization: Bearer <token>`)

## API Quick Examples

### 1. Create Staff

```bash
curl -X POST http://localhost:8088/staff/create \
  -H 'Content-Type: application/json' \
  -d '{
    "username":"alice",
    "password":"password123",
    "hospital":"HOSPITAL_A"
  }'
```

Response includes `token` for immediate use.

### 2. Staff Login

```bash
curl -X POST http://localhost:8088/staff/login \
  -H 'Content-Type: application/json' \
  -d '{
    "username":"alice",
    "password":"password123",
    "hospital":"HOSPITAL_A"
  }'
```

Save the returned `token` for subsequent requests.

### 3. Search Patients (Requires Auth)

Search all patients in your hospital:

```bash
TOKEN="<your_jwt_token>"
curl -X POST http://localhost:8088/patient/search \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{}'
```

Search by national ID:

```bash
curl -X POST http://localhost:8088/patient/search \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"national_id":"1234567890123"}'
```

Search by name (partial, case-insensitive):

```bash
curl -X POST http://localhost:8088/patient/search \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"first_name":"somchai"}'
```

All search fields are optional:
- `national_id` - exact match
- `passport_id` - exact match
- `first_name` - partial match (Thai or English)
- `middle_name` - partial match (Thai or English)
- `last_name` - partial match (Thai or English)
- `date_of_birth` - exact match (YYYY-MM-DD)
- `phone_number` - exact match
- `email` - case-insensitive exact match

- Implement staff APIs: `/staff/create`, `/staff/login`
- Implement secured patient search: `/patient/search`
- Add unit tests for all APIs
