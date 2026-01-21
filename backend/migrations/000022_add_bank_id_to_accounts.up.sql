-- Add bank_id column to accounts table (nullable first)
ALTER TABLE accounts
ADD COLUMN bank_id UUID;

-- Add foreign key constraint
ALTER TABLE accounts
ADD CONSTRAINT fk_accounts_bank FOREIGN KEY (bank_id) REFERENCES banks(id) ON DELETE RESTRICT;

-- Create index on bank_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_accounts_bank ON accounts(bank_id);

-- Note: bank_id will be set to NOT NULL after initial sync
-- This allows existing accounts to be updated with a bank after sync
