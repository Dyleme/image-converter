CREATE TABLE IF NOT EXISTS users (
  id             SERIAL UNIQUE PRIMARY KEY,
  nickname       VARCHAR(50) UNIQUE NOT NULL,
  email          VARCHAR(50) NOT NULL,
  password_hash  VARCHAR(250) NOT NULL
);

-- INSERT INTO users
--   (nickname, email, password_hash) 
-- VALUES 
--   ('Dylefme', 'Alex@', 'fdjklgj'),
--   ('Nick', 'Nickname@', 'hkldsf'),
--   ('Jerdsfu', 'Gerry Mulligan@', 'dsjlm'),
--   ('Sarasdefh Vaughan', 'Sarah Vaughan@', 'sdjfk');

CREATE TYPE operation_status AS ENUM ('queued', 'processing', 'done');

CREATE TYPE image_type AS ENUM ('jpeg', 'png');

CREATE TABLE IF NOT EXISTS requests (
  id                  SERIAL UNIQUE PRIMARY KEY,
  op_status           operation_status NOT NULL DEFAULT 'queued',
  request_time        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  completion_time     TIMESTAMP WITH TIME ZONE,
  original_id         INTEGER NOT NULL,
  processed_id        INTEGER,
  user_id             INTEGER NOT NULL,
  ratio               FLOAT NOT NULL DEFAULT 1,
  original_type       image_type NOT NULL,
  processed_type      image_type NOT NULL
);

CREATE TABLE IF NOT EXISTS images (
  id               SERIAL UNIQUE PRIMARY KEY,
  resoolution_x    INTEGER NOT NULL,
  resoolution_y    INTEGER NOT NULL,
  im_type          image_type NOT NULL,
  image_url        VARCHAR(50) NOT NULL,
  user_id          INTEGER NOT NULL,
  request_id       INTEGER 
);

-- INSERT INTO images(resoolution_x, resoolution_y, im_type, image_url, user_id, request_id)
-- VALUES (1080, 720, 'JPEG', 'image.url', 1, 1);

-- INSERT INTO requests(op_status, original_id, user_id, original_type, processed_type)
-- VALUES 
-- ('queued', 1, 1, 'JPEG', 'JPEG'),
-- ('done', 2, 1, 'PNG', 'JPEG'),
-- ('processing', 3, 2, 'PNG', 'JPEG');