-- +migrate Up
CREATE TABLE tokens(
  token TEXT NOT NULL UNIQUE,
  username TEXT NOT NULL,
  description TEXT
);

-- +migrate Down
DROP TABLE tokens;
