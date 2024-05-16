USE gosnip;

CREATE TABLE IF NOT EXISTS sessions (
 token CHAR(43) PRIMARY KEY,
 data BLOB NOT NULL,
 expiry TIMESTAMP(6) NOT NULL
);

CREATE PROCEDURE CreateSessionsExpiryIndexIfNotExists()
BEGIN
  DECLARE indexExists INTEGER;

  SELECT COUNT(1) INTO indexExists
  FROM INFORMATION_SCHEMA.STATISTICS
  WHERE table_schema = 'gosnip'
    AND table_name = 'sessions'
    AND index_name = 'sessions_expiry_idx';

  IF indexExists = 0 THEN
    CREATE INDEX sessions_expiry_idx ON sessions(expiry);
  END IF;
END;

CALL CreateSessionsExpiryIndexIfNotExists();

DROP PROCEDURE CreateSessionsExpiryIndexIfNotExists;