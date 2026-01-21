-- Create credit_card_bills table
CREATE TABLE IF NOT EXISTS credit_card_bills (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    credit_card_id UUID NOT NULL,
    month INTEGER NOT NULL,
    year INTEGER NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'open', -- 'open', 'closed', 'paid'
    closing_date DATE NOT NULL,
    due_date DATE NOT NULL,
    payment_transaction_id UUID NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_credit_card_bills_credit_card FOREIGN KEY (credit_card_id) REFERENCES credit_cards(id) ON DELETE CASCADE,
    CONSTRAINT fk_credit_card_bills_payment FOREIGN KEY (payment_transaction_id) REFERENCES transactions(id) ON DELETE SET NULL,
    CONSTRAINT check_bill_status CHECK (status IN ('open', 'closed', 'paid')),
    CONSTRAINT check_month CHECK (month >= 1 AND month <= 12),
    CONSTRAINT check_year CHECK (year >= 2000 AND year <= 2100),
    CONSTRAINT unique_bill_per_card_month UNIQUE (credit_card_id, year, month)
);

-- Create index on credit_card_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_credit_card_bills_credit_card ON credit_card_bills(credit_card_id);

-- Create index on (credit_card_id, year, month) for queries
CREATE INDEX IF NOT EXISTS idx_credit_card_bills_card_period ON credit_card_bills(credit_card_id, year, month);

-- Create index on payment_transaction_id
CREATE INDEX IF NOT EXISTS idx_credit_card_bills_payment ON credit_card_bills(payment_transaction_id) WHERE payment_transaction_id IS NOT NULL;

-- Create index on status for filtering
CREATE INDEX IF NOT EXISTS idx_credit_card_bills_status ON credit_card_bills(credit_card_id, status);

-- Create index on due_date for calendar queries
CREATE INDEX IF NOT EXISTS idx_credit_card_bills_due_date ON credit_card_bills(due_date) WHERE status IN ('closed', 'open');
