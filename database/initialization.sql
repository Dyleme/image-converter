CREATE TABLE IF NOT EXISTS users;
CREATE TABLE album (
  id         SERIAL PRIMARY KEY,
  nickname   VARCHAR(50) NOT NULL,
  email     VARCHAR(50) NOT NULL, 
  password_hash  DECIMAL(5,2) NOT NULL,
);

INSERT INTO users
  (nickname, email, password_hash) 
VALUES 
  ('Dyleme', 'Alex@', 'fdjklgj'),
  ('Nick', 'Nickname@', 'hkldsf'),
  ('Jeru', 'Gerry Mulligan@', 'dsjlm'),
  ('Sarah Vaughan', 'Sarah Vaughan@', 'sdjfk');