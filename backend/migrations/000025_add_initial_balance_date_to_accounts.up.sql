-- Add initial_balance_date and initial_balance columns to accounts table
ALTER TABLE accounts
ADD COLUMN IF NOT EXISTS initial_balance_date TIMESTAMP NULL,
ADD COLUMN IF NOT EXISTS initial_balance BIGINT NULL;

-- Add comment to explain the columns
COMMENT ON COLUMN accounts.initial_balance_date IS 'Data de referência do saldo inicial da conta';
COMMENT ON COLUMN accounts.initial_balance IS 'Saldo inicial da conta em centavos na data de referência';
