-- Switch to using the `gosnip` database.
USE gosnip;

-- Create a `snippets` table.
CREATE TABLE IF NOT EXISTS snippets (
 id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
 title VARCHAR(100) NOT NULL,
 content TEXT NOT NULL,
 created DATETIME NOT NULL,
 expires DATETIME NOT NULL
);
-- Add an index on the created column.
DELIMITER //
CREATE PROCEDURE CreateIndexIfNotExists()
BEGIN
  DECLARE indexExists INTEGER;

  SELECT COUNT(1) INTO indexExists
  FROM INFORMATION_SCHEMA.STATISTICS
  WHERE table_schema = 'gosnip'
    AND table_name = 'snippets'
    AND index_name = 'idx_snippets_created';

  IF indexExists = 0 THEN
    CREATE INDEX idx_snippets_created ON snippets(created);
  END IF;
END //
DELIMITER ;

CALL CreateIndexIfNotExists();

DROP PROCEDURE CreateIndexIfNotExists;
