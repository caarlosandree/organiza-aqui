-- Remove financial integration fields from tasks table
DROP INDEX IF EXISTS idx_tasks_financial_account;
ALTER TABLE tasks
DROP COLUMN IF EXISTS financial_category_id,
DROP COLUMN IF EXISTS financial_amount,
DROP COLUMN IF EXISTS financial_account_id;
