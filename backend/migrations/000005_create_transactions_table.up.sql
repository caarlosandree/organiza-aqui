-- Create transactions table
CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    account_id UUID NOT NULL,
    category_id UUID NULL,
    type VARCHAR(50) NOT NULL, -- 'income', 'expense', 'transfer'
    amount BIGINT NOT NULL, -- em centavos
    description TEXT,
    date DATE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_transactions_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_transactions_account FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE,
    CONSTRAINT fk_transactions_category FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL,
    CONSTRAINT check_transaction_amount CHECK (amount > 0)
);

-- Create index on user_id and date for faster lookups and sorting
CREATE INDEX IF NOT EXISTS idx_transactions_user_date ON transactions(user_id, date DESC);

-- Create index on account_id for account-specific queries
CREATE INDEX IF NOT EXISTS idx_transactions_account ON transactions(account_id);

-- Create index on category_id for category-specific queries
CREATE INDEX IF NOT EXISTS idx_transactions_category ON transactions(category_id);

-- Create index on type for filtering
CREATE INDEX IF NOT EXISTS idx_transactions_type ON transactions(user_id, type);
