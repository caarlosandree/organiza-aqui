-- Make bank_id NOT NULL after ensuring all accounts have a bank
-- This migration should be run after the initial bank sync
ALTER TABLE accounts
ALTER COLUMN bank_id SET NOT NULL;
