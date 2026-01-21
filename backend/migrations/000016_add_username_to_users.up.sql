-- Add username column to users table (nullable first)
ALTER TABLE users ADD COLUMN IF NOT EXISTS username VARCHAR(50);

-- Create index on username for faster lookups
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);

-- Update existing users with a default username based on email (if username is null)
-- This ensures existing users have a username
-- Generate unique usernames by appending numbers if needed
DO $$
DECLARE
    user_record RECORD;
    base_username TEXT;
    final_username TEXT;
    counter INTEGER;
BEGIN
    FOR user_record IN SELECT id, email FROM users WHERE username IS NULL LOOP
        base_username := LOWER(REGEXP_REPLACE(SPLIT_PART(user_record.email, '@', 1), '[^a-z0-9]', '', 'g'));
        final_username := base_username;
        counter := 1;
        
        -- Check if username already exists, if so, append number
        WHILE EXISTS (SELECT 1 FROM users WHERE username = final_username) LOOP
            final_username := base_username || counter::TEXT;
            counter := counter + 1;
        END LOOP;
        
        UPDATE users SET username = final_username WHERE id = user_record.id;
    END LOOP;
END $$;

-- Add unique constraint after populating all usernames
ALTER TABLE users ADD CONSTRAINT users_username_unique UNIQUE (username);

-- Make username NOT NULL after setting defaults
ALTER TABLE users ALTER COLUMN username SET NOT NULL;
