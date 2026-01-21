-- Create habits table
CREATE TABLE IF NOT EXISTS habits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    color VARCHAR(7) NOT NULL DEFAULT '#3b82f6',
    frequency VARCHAR(20) NOT NULL DEFAULT 'daily', -- 'daily', 'weekly', 'monthly'
    target_days INTEGER NOT NULL DEFAULT 1, -- Quantos dias por período (ex: 3 vezes por semana = 3)
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_habit_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT check_frequency CHECK (frequency IN ('daily', 'weekly', 'monthly'))
);

-- Create habit_tracking table (registro de execução dos hábitos)
CREATE TABLE IF NOT EXISTS habit_tracking (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    habit_id UUID NOT NULL,
    date DATE NOT NULL,
    completed BOOLEAN NOT NULL DEFAULT true,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_habit_tracking_habit FOREIGN KEY (habit_id) REFERENCES habits(id) ON DELETE CASCADE,
    CONSTRAINT unique_habit_date UNIQUE (habit_id, date)
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_habits_user ON habits(user_id);
CREATE INDEX IF NOT EXISTS idx_habit_tracking_habit ON habit_tracking(habit_id);
CREATE INDEX IF NOT EXISTS idx_habit_tracking_date ON habit_tracking(habit_id, date DESC);
