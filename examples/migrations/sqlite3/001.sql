-- +migrate Up
insert into todos(name, description) values("foo", "bar");
insert into todos(name, description) values("foo", "bar");
insert into todos(name, description) values("foo", "bar");
insert into todos(name, description) values("foo", "bar");
insert into todos(name, description) values("foo", "bar");
insert into todos(name, description) values("foo", "bar");
insert into todos(name, description) values("foo", "bar");
insert into todos(name, description) values("foo", "bar");
insert into todos(name, description) values("foo", "bar");

-- +migrate Down
truncate table todos;
