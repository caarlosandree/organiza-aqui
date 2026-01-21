-- Create calendar_events table
CREATE TABLE IF NOT EXISTS calendar_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP,
    all_day BOOLEAN NOT NULL DEFAULT false,
    location VARCHAR(255),
    color VARCHAR(7) NOT NULL DEFAULT '#3b82f6',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_calendar_event_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_calendar_events_user ON calendar_events(user_id);
CREATE INDEX IF NOT EXISTS idx_calendar_events_start_date ON calendar_events(user_id, start_date);
CREATE INDEX IF NOT EXISTS idx_calendar_events_date_range ON calendar_events(user_id, start_date, end_date);
