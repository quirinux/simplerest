-- +migrate Up
CREATE TABLE todos(
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT
);

-- +migrate Down
DROP TABLE todos;
