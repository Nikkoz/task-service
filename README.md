# Task Service

Task Service is a pet-project REST API written in Go.

The project demonstrates:
- layered architecture
- domain value objects
- PostgreSQL integration
- SQL migrations
- unit / integration / HTTP tests
- Docker-based local development

## Stack

- Go
- Gin
- PostgreSQL
- pgx
- golang-migrate
- mockery
- Docker Compose

## Project structure

```text
cmd/
  api/                  # application entrypoint

internal/
  app/                  # application bootstrap
  config/               # config loading
  domain/               # domain entities and value objects
  repository/           # repository interfaces and postgres implementation
  service/              # use cases / business logic
  transport/            # HTTP layer
  testutil/             # integration test helpers

migrations/             # database migrations
deployments/            # docker-compose and deployment-related files
tools/                  # pinned dev tools
```

## Data model

### PostgreSQL

This service uses PostgreSQL. Database schema is managed via migrations (`migrations/*.up.sql`, `migrations/*.down.sql`).

### Enum: `task_status`

`task_status` is a PostgreSQL ENUM:

- `planned`
- `in_progress`
- `done`

### Table: `tasks`

Stores tasks (v1: single-user scope).

Columns:

- `id` — `bigint generated always as identity`, primary key
- `title` — `varchar(255)`, not null
- `description` — `text`, not null, default `''`
- `status` — `task_status`, not null, default `'planned'`
- `due_date` — `timestamptz`, nullable
- `created_at` — `timestamptz`, not null, default `now()`
- `updated_at` — `timestamptz`, not null, default `now()`

Behavior:

- `updated_at` is automatically updated on every UPDATE via trigger `trg_tasks_set_updated_at`
  (trigger function: `trigger_set_updated_at()`).

Indexes:

- `idx_tasks_due_date` on `(due_date)`
- `idx_tasks_status_due_date` on `(status, due_date)` (covers filtering by status and sorting/range by due_date inside status)

# API

## Health
- `GET /health`

## Tasks
- `GET /tasks`
- `POST /tasks`
- `GET /tasks/{id}`
- `PUT /tasks/{id}`
- `DELETE /tasks/{id}`

## Example requests
### Create a new task

```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Buy milk",
    "description": "2 liters",
    "status": "planned",
    "due_date": null
  }'
```

### Update a task

```bash
curl -X PUT http://localhost:8080/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Buy oat milk",
    "description": "barista edition",
    "status": "done",
    "due_date": null
  }'
```

### Delete a task
```bash
curl -X DELETE http://localhost:8080/tasks/1
```

### Get task by ID
```bash
curl http://localhost:8080/tasks/1
```

# Configuration

Application configuration is loaded from environment variables.
For local development:
- `.env` — main local environment
- `.env.testing` — environment for integration tests

