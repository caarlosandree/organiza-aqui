-- Remove indexes
DROP INDEX IF EXISTS idx_transactions_installment;
DROP INDEX IF EXISTS idx_transactions_parent;
DROP INDEX IF EXISTS idx_transactions_to_account;
DROP INDEX IF EXISTS idx_transactions_status;
DROP INDEX IF EXISTS idx_transactions_external_id;

-- Remove constraints
ALTER TABLE transactions DROP CONSTRAINT IF EXISTS check_transaction_status;
ALTER TABLE transactions DROP CONSTRAINT IF EXISTS fk_transactions_parent;
ALTER TABLE transactions DROP CONSTRAINT IF EXISTS fk_transactions_to_account;

-- Remove columns
ALTER TABLE transactions 
DROP COLUMN IF EXISTS external_id,
DROP COLUMN IF EXISTS total_installments,
DROP COLUMN IF EXISTS installment_number,
DROP COLUMN IF EXISTS parent_transaction_id,
DROP COLUMN IF EXISTS to_account_id,
DROP COLUMN IF EXISTS tags,
DROP COLUMN IF EXISTS status;
