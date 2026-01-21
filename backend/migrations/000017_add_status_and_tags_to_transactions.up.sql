-- Add new columns to transactions table
ALTER TABLE transactions 
ADD COLUMN IF NOT EXISTS status VARCHAR(50) DEFAULT 'paid',
ADD COLUMN IF NOT EXISTS tags TEXT[] DEFAULT '{}',
ADD COLUMN IF NOT EXISTS to_account_id UUID NULL,
ADD COLUMN IF NOT EXISTS parent_transaction_id UUID NULL,
ADD COLUMN IF NOT EXISTS installment_number INTEGER NULL,
ADD COLUMN IF NOT EXISTS total_installments INTEGER NULL,
ADD COLUMN IF NOT EXISTS external_id VARCHAR(255) NULL;

-- Add foreign key constraints
ALTER TABLE transactions
ADD CONSTRAINT fk_transactions_to_account FOREIGN KEY (to_account_id) REFERENCES accounts(id) ON DELETE SET NULL,
ADD CONSTRAINT fk_transactions_parent FOREIGN KEY (parent_transaction_id) REFERENCES transactions(id) ON DELETE CASCADE;

-- Add check constraint for status
ALTER TABLE transactions
ADD CONSTRAINT check_transaction_status CHECK (status IN ('pending', 'paid', 'cancelled'));

-- Add unique constraint on external_id for deduplication
CREATE UNIQUE INDEX IF NOT EXISTS idx_transactions_external_id ON transactions(external_id) WHERE external_id IS NOT NULL;

-- Create index on status for filtering
CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions(user_id, status);

-- Create index on to_account_id for transfer queries
CREATE INDEX IF NOT EXISTS idx_transactions_to_account ON transactions(to_account_id) WHERE to_account_id IS NOT NULL;

-- Create index on parent_transaction_id for installment queries
CREATE INDEX IF NOT EXISTS idx_transactions_parent ON transactions(parent_transaction_id) WHERE parent_transaction_id IS NOT NULL;

-- Create index on installment fields for installment queries
CREATE INDEX IF NOT EXISTS idx_transactions_installment ON transactions(user_id, parent_transaction_id, installment_number) WHERE parent_transaction_id IS NOT NULL;

-- Update existing transactions to have 'paid' status
UPDATE transactions SET status = 'paid' WHERE status IS NULL;
