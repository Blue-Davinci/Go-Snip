-- Switch to using the `gosnip` database.
USE gosnip;

-- Drop the index on the created column.
DROP INDEX idx_snippets_created ON snippets;

-- Drop the `snippets` table.
DROP TABLE IF EXISTS snippets;