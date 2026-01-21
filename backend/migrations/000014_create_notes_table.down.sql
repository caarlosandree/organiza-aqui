-- Drop notes table
DROP INDEX IF EXISTS idx_notes_tags;
DROP INDEX IF EXISTS idx_notes_pinned;
DROP INDEX IF EXISTS idx_notes_created_at;
DROP INDEX IF EXISTS idx_notes_user;
DROP TABLE IF EXISTS notes;
