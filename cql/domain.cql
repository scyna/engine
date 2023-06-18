USE scyna2;

/*module*/
DROP TABLE IF EXISTS module;
CREATE TABLE module
(
    code        text PRIMARY KEY,
    secret      text,
);

/*client*/
DROP TABLE IF EXISTS client;
CREATE TABLE client
(
    id          text PRIMARY KEY,
    secret      text, 
);

DROP TABLE IF EXISTS  client_use_endpoint;
CREATE TABLE  client_use_endpoint
(
    client      text,   
    url         text,   
    PRIMARY KEY (client, url)
);

/*setting*/
DROP TABLE IF EXISTS  setting;
CREATE TABLE  setting
(
    module     text,
    key         text,
    value       text,
    PRIMARY KEY (module, key)
);
