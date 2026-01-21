-- Drop additional financial indexes
DROP INDEX IF EXISTS idx_categories_user_path;
DROP INDEX IF EXISTS idx_transactions_account_date;
DROP INDEX IF EXISTS idx_transactions_user_date_type;
