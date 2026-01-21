-- Drop calendar_events table
DROP INDEX IF EXISTS idx_calendar_events_date_range;
DROP INDEX IF EXISTS idx_calendar_events_start_date;
DROP INDEX IF EXISTS idx_calendar_events_user;
DROP TABLE IF EXISTS calendar_events;
