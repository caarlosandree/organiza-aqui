-- Remove initial_balance_date and initial_balance columns from accounts table
ALTER TABLE accounts
DROP COLUMN IF EXISTS initial_balance_date,
DROP COLUMN IF EXISTS initial_balance;
