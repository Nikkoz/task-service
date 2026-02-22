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