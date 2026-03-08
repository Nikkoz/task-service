-- 003_create_users.up.sql

-- Table
CREATE TABLE users (
    id            BIGINT GENERATED ALWAYS AS IDENTITY,
    email         VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(250) NOT NULL,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT now(),

    CONSTRAINT pk_users PRIMARY KEY (id)
);

-- Trigger for updated_at
CREATE TRIGGER trg_users_set_updated_at
    BEFORE UPDATE
    ON users
    FOR EACH ROW
EXECUTE FUNCTION trigger_set_updated_at();

-- Indexes
CREATE UNIQUE INDEX uq_users_email_lower
    ON users (lower(email));

-- Comments
COMMENT ON TABLE users IS 'Application users';
COMMENT ON COLUMN users.id IS 'Primary key';
COMMENT ON COLUMN users.email IS 'User email, unique case-insensitive';
COMMENT ON COLUMN users.password_hash IS 'Password hash';
COMMENT ON COLUMN users.created_at IS 'Creation timestamp';
COMMENT ON COLUMN users.updated_at IS 'Last update timestamp';

