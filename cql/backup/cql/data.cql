/*core context*/
insert into scyna.module(code,secret) values ( 'scyna_engine', '123456');
insert into scyna.module(code,secret) values ( 'scyna_manager', '123456');
insert into scyna.module(code,secret) values ( 'scyna_test', '123456');

insert into scyna.client(id, secret) VALUES ('admin', 'admin');
insert into scyna.application(code, auth_url) VALUES ('noauth', '/scyna/noauth');
