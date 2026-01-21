-- Drop indexes
DROP INDEX IF EXISTS idx_credit_cards_created_at;
DROP INDEX IF EXISTS idx_credit_cards_account;
DROP INDEX IF EXISTS idx_credit_cards_user;

-- Drop table
DROP TABLE IF EXISTS credit_cards;
