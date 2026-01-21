-- Create accounts table
CREATE TABLE IF NOT EXISTS accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL, -- 'checking', 'savings', 'credit', 'investment'
    balance BIGINT NOT NULL DEFAULT 0, -- em centavos
    currency VARCHAR(3) NOT NULL DEFAULT 'BRL',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_accounts_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT unique_account_name_per_user UNIQUE (user_id, name)
);

-- Create index on user_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_accounts_user ON accounts(user_id);

-- Create index on created_at for sorting
CREATE INDEX IF NOT EXISTS idx_accounts_created_at ON accounts(created_at DESC);
