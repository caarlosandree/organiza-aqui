-- Drop timeline_events table
DROP INDEX IF EXISTS idx_timeline_type_date;
DROP INDEX IF EXISTS idx_timeline_entity;
DROP INDEX IF EXISTS idx_timeline_date;
DROP INDEX IF EXISTS idx_timeline_user;
DROP TABLE IF EXISTS timeline_events;
