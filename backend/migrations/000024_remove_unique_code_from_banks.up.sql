-- Remove UNIQUE constraint from code column
-- Some banks can have the same code but different ISPB
ALTER TABLE banks DROP CONSTRAINT IF EXISTS banks_code_key;
