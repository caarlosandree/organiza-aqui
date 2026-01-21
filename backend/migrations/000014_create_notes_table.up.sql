-- Create notes table
CREATE TABLE IF NOT EXISTS notes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    tags TEXT[], -- Array de tags
    is_pinned BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_note_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_notes_user ON notes(user_id);
CREATE INDEX IF NOT EXISTS idx_notes_created_at ON notes(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_notes_pinned ON notes(user_id, is_pinned, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_notes_tags ON notes USING GIN(tags); -- GIN index para busca em arrays
