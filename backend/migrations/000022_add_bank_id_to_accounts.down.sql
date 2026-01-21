-- Remove bank_id column from accounts table
DROP INDEX IF EXISTS idx_accounts_bank;
ALTER TABLE accounts DROP CONSTRAINT IF EXISTS fk_accounts_bank;
ALTER TABLE accounts DROP COLUMN IF EXISTS bank_id;
