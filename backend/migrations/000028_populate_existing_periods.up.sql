-- Migration to populate periods for existing transactions
-- This creates periods based on transaction date and account type, then updates period_id

DO $$
DECLARE
    trans_record RECORD;
    period_record RECORD;
    period_id_val UUID;
    period_type_val VARCHAR(50);
    reference_date DATE;
BEGIN
    -- Loop through all existing transactions
    FOR trans_record IN 
        SELECT t.id, t.user_id, t.account_id, t.date, a.type as account_type
        FROM transactions t
        INNER JOIN accounts a ON t.account_id = a.id
        WHERE t.period_id IS NULL
    LOOP
        -- Determine period_type based on account type
        IF trans_record.account_type = 'credit' THEN
            period_type_val := 'credit_card';
        ELSE
            period_type_val := 'bank';
        END IF;
        
        -- Use transaction date as reference_month (first day of the month)
        reference_date := DATE_TRUNC('month', trans_record.date)::DATE;
        
        -- Check if period already exists
        SELECT id INTO period_record
        FROM transaction_periods
        WHERE account_id = trans_record.account_id
          AND period_type = period_type_val
          AND year = EXTRACT(YEAR FROM reference_date)::INTEGER
          AND month = EXTRACT(MONTH FROM reference_date)::INTEGER;
        
        -- If period doesn't exist, create it
        IF period_record IS NULL THEN
            INSERT INTO transaction_periods (id, user_id, account_id, period_type, year, month, status, created_at, updated_at)
            VALUES (
                gen_random_uuid(),
                trans_record.user_id,
                trans_record.account_id,
                period_type_val,
                EXTRACT(YEAR FROM reference_date)::INTEGER,
                EXTRACT(MONTH FROM reference_date)::INTEGER,
                'open',
                NOW(),
                NOW()
            )
            RETURNING id INTO period_id_val;
        ELSE
            period_id_val := period_record.id;
        END IF;
        
        -- Update transaction with period_id and reference_month
        UPDATE transactions
        SET period_id = period_id_val,
            reference_month = reference_date
        WHERE id = trans_record.id;
    END LOOP;
END $$;
