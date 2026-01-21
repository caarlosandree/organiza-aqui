-- Create recurrence_patterns table
CREATE TABLE IF NOT EXISTS recurrence_patterns (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    transaction_id UUID NOT NULL,
    frequency VARCHAR(50) NOT NULL, -- 'daily', 'weekly', 'monthly', 'yearly'
    interval INTEGER NOT NULL DEFAULT 1, -- a cada X dias/semanas/meses
    end_date DATE NULL, -- NULL = infinito
    last_generated_date DATE NULL, -- última data em que foi gerada uma transação
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_recurrence_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_recurrence_transaction FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE CASCADE,
    CONSTRAINT check_frequency CHECK (frequency IN ('daily', 'weekly', 'monthly', 'yearly')),
    CONSTRAINT check_interval CHECK (interval > 0)
);

-- Create index on user_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_recurrence_user ON recurrence_patterns(user_id);

-- Create index on transaction_id
CREATE INDEX IF NOT EXISTS idx_recurrence_transaction ON recurrence_patterns(transaction_id);

-- Create index on last_generated_date for queries de geração
CREATE INDEX IF NOT EXISTS idx_recurrence_last_generated ON recurrence_patterns(last_generated_date);
