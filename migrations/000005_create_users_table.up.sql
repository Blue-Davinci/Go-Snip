USE gosnip;

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
  WHERE CONSTRAINT_SCHEMA = 'gosnip'
    AND TABLE_NAME = 'users'
    AND CONSTRAINT_NAME = 'users_uc_email';

  IF constraintExists = 0 THEN
    ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);
  END IF;
END;

CALL AddConstraintIfNotExists();

DROP PROCEDURE AddConstraintIfNotExists;