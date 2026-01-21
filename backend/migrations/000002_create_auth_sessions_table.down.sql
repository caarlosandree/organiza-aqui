-- Drop auth_sessions table
DROP INDEX IF EXISTS idx_auth_sessions_token;
DROP INDEX IF EXISTS idx_auth_sessions_user_expires;
DROP TABLE IF EXISTS auth_sessions;
