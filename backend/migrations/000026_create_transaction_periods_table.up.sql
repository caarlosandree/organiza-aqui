-- Create transaction_periods table
CREATE TABLE IF NOT EXISTS transaction_periods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    account_id UUID NOT NULL,
    period_type VARCHAR(50) NOT NULL, -- 'bank' ou 'credit_card'
    year INTEGER NOT NULL,
    month INTEGER NOT NULL, -- 1-12
    status VARCHAR(50) NOT NULL DEFAULT 'open', -- 'open', 'closed', 'archived'
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_transaction_periods_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_transaction_periods_account FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE,
    CONSTRAINT check_period_type CHECK (period_type IN ('bank', 'credit_card')),
    CONSTRAINT check_status CHECK (status IN ('open', 'closed', 'archived')),
    CONSTRAINT check_month CHECK (month >= 1 AND month <= 12),
    CONSTRAINT check_year CHECK (year >= 2000 AND year <= 2100),
    CONSTRAINT unique_period_per_account UNIQUE (account_id, period_type, year, month)
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_transaction_periods_user ON transaction_periods(user_id);
CREATE INDEX IF NOT EXISTS idx_transaction_periods_account ON transaction_periods(account_id);
CREATE INDEX IF NOT EXISTS idx_transaction_periods_period ON transaction_periods(account_id, period_type, year, month);
CREATE INDEX IF NOT EXISTS idx_transaction_periods_year_month ON transaction_periods(year, month DESC);
