-- 004_add_user_id_to_tasks.down.sql

ALTER TABLE tasks DROP CONSTRAINT fk_tasks_user;
ALTER TABLE tasks DROP COLUMN user_id;