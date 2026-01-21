-- Revert bank_id to nullable
ALTER TABLE accounts
ALTER COLUMN bank_id DROP NOT NULL;
