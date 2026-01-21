-- Create categories table
CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    parent_id UUID NULL,
    path VARCHAR(255) NOT NULL, -- Materialized Path (ex: "1.2.5")
    type VARCHAR(50) NOT NULL, -- 'income', 'expense'
    color VARCHAR(7) NOT NULL DEFAULT '#3b82f6', -- hex color
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_categories_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_categories_parent FOREIGN KEY (parent_id) REFERENCES categories(id) ON DELETE CASCADE,
    CONSTRAINT unique_category_name_per_user UNIQUE (user_id, name)
);

-- Create index on user_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_categories_user ON categories(user_id);

-- Create index on parent_id for hierarchical queries
CREATE INDEX IF NOT EXISTS idx_categories_parent ON categories(parent_id);

-- Create index on path for path-based queries
CREATE INDEX IF NOT EXISTS idx_categories_path ON categories(path);

-- Create index on type for filtering
CREATE INDEX IF NOT EXISTS idx_categories_type ON categories(user_id, type);
