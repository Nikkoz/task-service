-- 001_base_triggers.up.sql

-- Universal trigger function: sets updated_at = now() on UPDATE
-- Works for any table that has an "updated_at" column.
CREATE OR REPLACE FUNCTION trigger_set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = now();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION trigger_set_updated_at() IS
  'Universal BEFORE UPDATE trigger function to set updated_at to now()';