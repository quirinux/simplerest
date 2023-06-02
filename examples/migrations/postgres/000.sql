-- +migrate Up
CREATE TABLE lists(
  id INT GENERATED ALWAYS AS IDENTITY,
  name TEXT NOT NULL,
  description TEXT,
  PRIMARY KEY(id)
);

-- +migrate Down
DROP TABLE lists;
