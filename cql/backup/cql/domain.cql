/*module*/
DROP TABLE IF EXISTS scyna.module;
CREATE TABLE scyna.module
(
    code        text PRIMARY KEY,
    secret      text,
);

/*client*/
DROP TABLE IF EXISTS scyna.client;
CREATE TABLE scyna.client
(
    id          text PRIMARY KEY,
    secret      text, 
);

/*application*/
DROP TABLE IF EXISTS scyna.application;
CREATE TABLE scyna.application
(
    code        text PRIMARY KEY,
    auth_url    text,
);

/*relations*/
DROP TABLE IF EXISTS scyna.client_use_endpoint;
CREATE TABLE scyna.client_use_endpoint
(
    client      text,   /*client id*/  
    url         text,   /*endpoint path*/
    PRIMARY KEY (client, url)
);

DROP TABLE IF EXISTS scyna.application_use_endpoint;
CREATE TABLE scyna.application_use_endpoint
(
    application     text,   /*application code*/  
    url             text,   /*endpoint path*/
    PRIMARY KEY (application, url)
);

DROP TABLE IF EXISTS scyna.authentication;
CREATE TABLE scyna.authentication
(
    id      text PRIMARY KEY,
    time    timestamp,
    expired timestamp,
    apps    list<text>,
    uid     text,
);

DROP TABLE IF EXISTS scyna.app_has_auth;
CREATE TABLE scyna.app_has_auth
(
    app     text,
    auth    text,
    uid     text,
    PRIMARY KEY (app, auth)
);


/*setting*/
DROP TABLE IF EXISTS scyna.setting;
CREATE TABLE scyna.setting
(
    module     text,
    key         text,
    value       text,
    PRIMARY KEY (module, key)
);