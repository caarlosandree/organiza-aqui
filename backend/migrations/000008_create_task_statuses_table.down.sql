-- Drop task_statuses table
DROP INDEX IF EXISTS idx_task_statuses_order;
DROP INDEX IF EXISTS idx_task_statuses_user;
DROP TABLE IF EXISTS task_statuses;
