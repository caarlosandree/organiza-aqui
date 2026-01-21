-- Add financial integration fields to tasks table
ALTER TABLE tasks
ADD COLUMN IF NOT EXISTS financial_account_id UUID REFERENCES accounts(id) ON DELETE SET NULL,
ADD COLUMN IF NOT EXISTS financial_amount BIGINT,
ADD COLUMN IF NOT EXISTS financial_category_id UUID REFERENCES categories(id) ON DELETE SET NULL;

-- Create index for financial queries
CREATE INDEX IF NOT EXISTS idx_tasks_financial_account ON tasks(financial_account_id);
