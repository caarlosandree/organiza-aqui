-- Create auth_sessions table
CREATE TABLE IF NOT EXISTS auth_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_auth_sessions_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create index on user_id and expires_at for faster lookups
CREATE INDEX IF NOT EXISTS idx_auth_sessions_user_expires ON auth_sessions(user_id, expires_at);

-- Create index on token_hash for token validation
CREATE INDEX IF NOT EXISTS idx_auth_sessions_token ON auth_sessions(token_hash);
