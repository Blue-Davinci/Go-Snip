-- Use our test DB
USE test_gosnip;

-- Create mock snippet table + index
CREATE TABLE IF NOT EXISTS snippets (
 id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
 title VARCHAR(100) NOT NULL,
 content TEXT NOT NULL,
 created DATETIME NOT NULL,
 expires DATETIME NOT NULL
);

CREATE PROCEDURE CreateIndexIfNotExists()
BEGIN
  DECLARE indexExists INTEGER;

  SELECT COUNT(1) INTO indexExists
  FROM INFORMATION_SCHEMA.STATISTICS
  WHERE table_schema = 'test_gosnip'
    AND table_name = 'snippets'
    AND index_name = 'idx_snippets_created';

  IF indexExists = 0 THEN
    CREATE INDEX idx_snippets_created ON snippets(created);    
  END IF;
END;

CALL CreateIndexIfNotExists();

DROP PROCEDURE CreateIndexIfNotExists;

-- Create mock user table + constraint
CREATE TABLE IF NOT EXISTS users (
 id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
 name VARCHAR(255) NOT NULL,
 email VARCHAR(255) NOT NULL,
 hashed_password CHAR(60) NOT NULL,
 created DATETIME NOT NULL
);

CREATE PROCEDURE AddConstraintIfNotExists()
BEGIN
  DECLARE constraintExists INTEGER;

  SELECT COUNT(1) INTO constraintExists
  FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS
  WHERE CONSTRAINT_SCHEMA = 'test_gosnip'
    AND TABLE_NAME = 'users'
    AND CONSTRAINT_NAME = 'users_uc_email';

  IF constraintExists = 0 THEN
    ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);
  END IF;
END;

CALL AddConstraintIfNotExists();

DROP PROCEDURE AddConstraintIfNotExists;

-- Add mock data
INSERT INTO users (name, email, hashed_password, created) VALUES (
 'Rob bor',
 'borrob@example.com',
 '$2a$12$NuTjWXm3KKntReFwyBVHyuf/to.HEwTy.eS206TNfkGfr6HzGJSWG',
 '2022-01-01 10:00:00'
);