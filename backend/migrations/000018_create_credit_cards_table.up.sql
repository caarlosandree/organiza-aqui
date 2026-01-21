-- Create credit_cards table
CREATE TABLE IF NOT EXISTS credit_cards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    account_id UUID NOT NULL,
    limit_amount BIGINT NOT NULL, -- em centavos
    closing_day INTEGER NOT NULL, -- dia do fechamento (1-31)
    due_day INTEGER NOT NULL, -- dia do vencimento (1-31)
    color VARCHAR(7) NOT NULL DEFAULT '#3b82f6',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_credit_cards_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_credit_cards_account FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE,
    CONSTRAINT check_closing_day CHECK (closing_day >= 1 AND closing_day <= 31),
    CONSTRAINT check_due_day CHECK (due_day >= 1 AND due_day <= 31),
    CONSTRAINT unique_credit_card_name_per_user UNIQUE (user_id, name)
);

-- Create index on user_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_credit_cards_user ON credit_cards(user_id);

-- Create index on account_id
CREATE INDEX IF NOT EXISTS idx_credit_cards_account ON credit_cards(account_id);

-- Create index on created_at for sorting
CREATE INDEX IF NOT EXISTS idx_credit_cards_created_at ON credit_cards(created_at DESC);
