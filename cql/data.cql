USE scyna2;

insert into module(code,secret,name) values ( 'scyna_engine', '123456', '');
insert into module(code,secret,name) values ( 'scyna_manager', '123456', '');
insert into module(code,secret,name) values ( 'scyna_test', '123456', '');
insert into client(id, secret) VALUES ('admin', 'admin');
insert into application(code, auth_url) VALUES ('noauth', '/scyna/noauth');
