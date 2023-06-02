-- +migrate Up
insert into tokens(token, username, description) values("467aa100-7883-4cbd-8152-b3478a0c3d0d", "foobar", "Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium.");

-- +migrate Down
truncate table tokens;
