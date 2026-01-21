-- Drop tasks table
DROP INDEX IF EXISTS idx_tasks_completed;
DROP INDEX IF EXISTS idx_tasks_lexorank;
DROP INDEX IF EXISTS idx_tasks_due_date;
DROP INDEX IF EXISTS idx_tasks_status;
DROP INDEX IF EXISTS idx_tasks_user;
DROP TABLE IF EXISTS tasks;
