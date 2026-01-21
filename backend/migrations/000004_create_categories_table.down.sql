-- Drop categories table
DROP INDEX IF EXISTS idx_categories_type;
DROP INDEX IF EXISTS idx_categories_path;
DROP INDEX IF EXISTS idx_categories_parent;
DROP INDEX IF EXISTS idx_categories_user;
DROP TABLE IF EXISTS categories;
