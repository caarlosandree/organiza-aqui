-- Reverse migration: Remove period_id and reference_month from existing transactions
-- Note: This does not delete the periods, only removes the references

UPDATE transactions
SET period_id = NULL,
    reference_month = NULL
WHERE period_id IS NOT NULL;
