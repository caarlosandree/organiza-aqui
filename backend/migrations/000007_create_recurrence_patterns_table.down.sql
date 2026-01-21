-- Drop recurrence_patterns table
DROP INDEX IF EXISTS idx_recurrence_last_generated;
DROP INDEX IF EXISTS idx_recurrence_transaction;
DROP INDEX IF EXISTS idx_recurrence_user;
DROP TABLE IF EXISTS recurrence_patterns;
