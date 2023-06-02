-- +migrate Up
-- +migrate StatementBegin
create or replace procedure transfer(
   todoid int,
   listid int
)
language plpgsql    
as $$
begin
    update todos 
    set list_id = listid 
    where id = todoid;

    commit;
end;$$
-- +migrate StatementEnd

-- +migrate Down
drop procedure transfer;
