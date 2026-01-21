-- Create task_statuses table
CREATE TABLE IF NOT EXISTS task_statuses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    color VARCHAR(7) NOT NULL DEFAULT '#000000',
    order_index INTEGER NOT NULL DEFAULT 0,
    is_default BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_task_status_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT unique_user_status_name UNIQUE (user_id, name)
);

-- Create index on user_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_task_statuses_user ON task_statuses(user_id);

-- Create index on order_index for sorting
CREATE INDEX IF NOT EXISTS idx_task_statuses_order ON task_statuses(user_id, order_index);
