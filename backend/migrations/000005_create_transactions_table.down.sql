-- Drop transactions table
DROP INDEX IF EXISTS idx_transactions_type;
DROP INDEX IF EXISTS idx_transactions_category;
DROP INDEX IF EXISTS idx_transactions_account;
DROP INDEX IF EXISTS idx_transactions_user_date;
DROP TABLE IF EXISTS transactions;
