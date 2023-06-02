-- +migrate Up
CREATE TABLE todos(
  id INT GENERATED ALWAYS AS IDENTITY,
  list_id INT,
  name TEXT NOT NULL,
  description TEXT,
  PRIMARY KEY(id),
  CONSTRAINT fk_list
    FOREIGN KEY(list_id)
      REFERENCES lists(id)
);

-- +migrate Down
DROP TABLE todos;
