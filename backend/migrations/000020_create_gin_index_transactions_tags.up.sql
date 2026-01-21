-- Create GIN index on tags array for efficient searching
CREATE INDEX IF NOT EXISTS idx_transactions_tags_gin ON transactions USING GIN(tags);
