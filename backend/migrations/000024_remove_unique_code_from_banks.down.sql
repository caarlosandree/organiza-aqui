-- Re-add UNIQUE constraint to code column
ALTER TABLE banks ADD CONSTRAINT banks_code_key UNIQUE (code);
