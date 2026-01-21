-- Drop accounts table
DROP INDEX IF EXISTS idx_accounts_created_at;
DROP INDEX IF EXISTS idx_accounts_user;
DROP TABLE IF EXISTS accounts;
