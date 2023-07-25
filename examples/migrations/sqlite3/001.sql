-- +migrate Up
insert into todos(name, description) values("name_foo", "desc bar");
insert into todos(name, description) values("name_foo", "desc bar");
insert into todos(name, description) values("name_foo", "desc bar");
insert into todos(name, description) values("name_foo", "desc bar");
insert into todos(name, description) values("name_foo", "desc bar");
insert into todos(name, description) values("name_foo", "desc bar");
insert into todos(name, description) values("name_foo", "desc bar");
insert into todos(name, description) values("name_foo", "desc bar");
insert into todos(name, description) values("name_foo", "desc bar");

-- +migrate Down
truncate table todos;
