USE scyna2;

/*module*/
DROP TABLE IF EXISTS module;
CREATE TABLE module
(
    code        text PRIMARY KEY,
    secret      text,
    name        text,
);

/*client*/
DROP TABLE IF EXISTS client;
CREATE TABLE client
(
    id          text PRIMARY KEY,
    secret      text,
    name        text, 
);

/*application*/
DROP TABLE IF EXISTS application;
CREATE TABLE  application
(
    code        text PRIMARY KEY,
    auth_url    text,
    name        text,
);

/*relations*/
DROP TABLE IF EXISTS  client_use_endpoint;
CREATE TABLE  client_use_endpoint
(
    client      text,   
    url         text,   
    PRIMARY KEY (client, url)
);

DROP TABLE IF EXISTS endpoint;
CREATE TABLE  endpoint
(
    url         text PRIMARY KEY,
    method      text,
    module      text,
    description text,
);

DROP TABLE IF EXISTS module_has_endpoint;
CREATE TABLE  module_has_endpoint
(
    module      text,
    url         text,
    PRIMARY KEY (module, url)
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
