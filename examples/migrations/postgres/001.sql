-- +migrate Up
insert into lists(name, description) values ('supermarket', 'supermaket shoppings list');
insert into lists(name, description) values ('shopping', 'shoppings list');
insert into lists(name, description) values ('house', 'house applinace fix list');

-- +migrate Down
truncate table lists;
