-- +migrate Up
CREATE TABLE todos(
  id NUMBER GENERATED BY DEFAULT AS IDENTITY,
  name VARCHAR(500) NOT NULL,
  description VARCHAR(500),
  PRIMARY KEY(id)
);

-- +migrate Down
DROP TABLE todos;
