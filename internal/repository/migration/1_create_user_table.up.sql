CREATE TABLE IF NOT EXISTS users (
  id             SERIAL UNIQUE PRIMARY KEY,
  nickname       VARCHAR(50) UNIQUE NOT NULL,
  email          VARCHAR(50) NOT NULL,
  password_hash  VARCHAR(250) NOT NULL
);