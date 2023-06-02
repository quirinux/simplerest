-- +migrate Up
-- +migrate StatementBegin
create or replace procedure raise_me()
language plpgsql    
as $$
DECLARE
  personal_email varchar(100) := 'raise@exceptions.com';
begin
  RAISE EXCEPTION 'Enter email is duplicate: %', personal_email
  USING HINT = 'Check email and enter correct email ID of user';
end;$$
-- +migrate StatementEnd

-- +migrate Down
drop procedure raise_me;
