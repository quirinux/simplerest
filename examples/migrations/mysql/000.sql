-- +migrate Up
CREATE TABLE todos(
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT
);

-- +migrate Down
DROP TABLE todos;
