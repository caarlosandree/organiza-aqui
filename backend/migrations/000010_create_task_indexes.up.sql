-- Additional indexes for task queries
CREATE INDEX IF NOT EXISTS idx_tasks_user_status_lexorank ON tasks(user_id, status_id, lexorank);
CREATE INDEX IF NOT EXISTS idx_tasks_user_priority ON tasks(user_id, priority);
CREATE INDEX IF NOT EXISTS idx_tasks_user_due_date_status ON tasks(user_id, due_date, status_id) WHERE due_date IS NOT NULL;
