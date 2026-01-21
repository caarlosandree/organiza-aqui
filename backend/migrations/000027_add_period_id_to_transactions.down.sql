-- Remove period_id and reference_month columns from transactions table
ALTER TABLE transactions DROP CONSTRAINT IF EXISTS fk_transactions_period;
DROP INDEX IF EXISTS idx_transactions_period;
DROP INDEX IF EXISTS idx_transactions_reference_month;
ALTER TABLE transactions DROP COLUMN IF EXISTS period_id;
ALTER TABLE transactions DROP COLUMN IF EXISTS reference_month;
