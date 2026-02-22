-- 002_create_tasks.up.sql

-- Enum type for task statuses
CREATE TYPE task_status AS ENUM ('planned', 'in_progress', 'done');

COMMENT ON TYPE task_status IS 'Task lifecycle status';

-- Table
CREATE TABLE tasks (
    id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title       VARCHAR(100) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    status      task_status NOT NULL DEFAULT 'planned',
    due_date    TIMESTAMPTZ NULL,

    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Trigger for updated_at
CREATE TRIGGER trg_tasks_set_updated_at
    BEFORE UPDATE ON tasks
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_updated_at();

-- Indexes
CREATE INDEX idx_tasks_due_date ON tasks (due_date);
CREATE INDEX idx_tasks_status_due_date ON tasks (status, due_date);

-- Comments
COMMENT ON TABLE tasks IS 'Tasks';
COMMENT ON COLUMN tasks.id IS 'Primary key';
COMMENT ON COLUMN tasks.title IS 'Title';
COMMENT ON COLUMN tasks.description IS 'Description';
COMMENT ON COLUMN tasks.status IS 'Status enum: planned | in_progress | done';
COMMENT ON COLUMN tasks.due_date IS 'Optional deadline';
COMMENT ON COLUMN tasks.created_at IS 'Creation timestamp';
COMMENT ON COLUMN tasks.updated_at IS 'Last update timestamp';