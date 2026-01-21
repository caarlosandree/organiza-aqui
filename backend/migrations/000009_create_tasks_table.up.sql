-- Create tasks table
CREATE TABLE IF NOT EXISTS tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    status_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    priority VARCHAR(20) NOT NULL DEFAULT 'medium', -- 'low', 'medium', 'high', 'urgent'
    due_date DATE,
    completed_at TIMESTAMP,
    lexorank VARCHAR(255) NOT NULL DEFAULT '0|', -- Lexorank para ordenação drag-and-drop
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_task_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_task_status FOREIGN KEY (status_id) REFERENCES task_statuses(id) ON DELETE RESTRICT,
    CONSTRAINT check_priority CHECK (priority IN ('low', 'medium', 'high', 'urgent'))
);

-- Create index on user_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_tasks_user ON tasks(user_id);

-- Create index on status_id
CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status_id);

-- Create index on due_date for filtering
CREATE INDEX IF NOT EXISTS idx_tasks_due_date ON tasks(user_id, due_date);

-- Create index on lexorank for sorting
CREATE INDEX IF NOT EXISTS idx_tasks_lexorank ON tasks(status_id, lexorank);

-- Create index on completed_at for filtering completed tasks
CREATE INDEX IF NOT EXISTS idx_tasks_completed ON tasks(user_id, completed_at);
