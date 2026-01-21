-- Remove index on username
DROP INDEX IF EXISTS idx_users_username;

-- Remove username column
ALTER TABLE users DROP COLUMN IF EXISTS username;
