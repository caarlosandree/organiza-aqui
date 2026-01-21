-- Drop indexes
DROP INDEX IF EXISTS idx_credit_card_bills_due_date;
DROP INDEX IF EXISTS idx_credit_card_bills_status;
DROP INDEX IF EXISTS idx_credit_card_bills_payment;
DROP INDEX IF EXISTS idx_credit_card_bills_card_period;
DROP INDEX IF EXISTS idx_credit_card_bills_credit_card;

-- Drop table
DROP TABLE IF EXISTS credit_card_bills;
