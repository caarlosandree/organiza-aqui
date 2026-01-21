-- Add period_id and reference_month columns to transactions table
ALTER TABLE transactions 
ADD COLUMN IF NOT EXISTS period_id UUID NULL,
ADD COLUMN IF NOT EXISTS reference_month DATE NULL;

-- Add foreign key constraint for period_id
ALTER TABLE transactions
ADD CONSTRAINT fk_transactions_period FOREIGN KEY (period_id) REFERENCES transaction_periods(id) ON DELETE SET NULL;

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_transactions_period ON transactions(period_id);
CREATE INDEX IF NOT EXISTS idx_transactions_reference_month ON transactions(reference_month);

-- Add comment to explain the columns
COMMENT ON COLUMN transactions.period_id IS 'ID do período de referência da transação';
COMMENT ON COLUMN transactions.reference_month IS 'Mês de referência da transação (primeiro dia do mês: YYYY-MM-01)';
