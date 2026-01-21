-- Additional indexes for financial module performance

-- Composite index for common transaction queries (user, date range, type)
CREATE INDEX IF NOT EXISTS idx_transactions_user_date_type ON transactions(user_id, date DESC, type);

-- Index for account balance calculations
CREATE INDEX IF NOT EXISTS idx_transactions_account_date ON transactions(account_id, date DESC);

-- Index for category hierarchy queries
CREATE INDEX IF NOT EXISTS idx_categories_user_path ON categories(user_id, path);
