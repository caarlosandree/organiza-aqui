-- Create banks table
CREATE TABLE IF NOT EXISTS banks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ispb VARCHAR(8) NOT NULL UNIQUE,
    code INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for faster lookups
CREATE INDEX IF NOT EXISTS idx_banks_ispb ON banks(ispb);
CREATE INDEX IF NOT EXISTS idx_banks_code ON banks(code);
CREATE INDEX IF NOT EXISTS idx_banks_name ON banks(name);

-- Note: code não é UNIQUE porque alguns bancos podem ter o mesmo código
-- mas ISPB diferente (ex: diferentes entidades da B3, CIP, etc.)
