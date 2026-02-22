-- 002_create_tasks.down.sql

DROP TRIGGER IF EXISTS trg_tasks_set_updated_at ON tasks;

DROP TABLE IF EXISTS tasks;

DROP TYPE IF EXISTS task_status;