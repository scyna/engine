/*domain*/
DROP TABLE IF EXISTS scyna.domain;
CREATE TABLE scyna.domain
(
    code     text PRIMARY KEY, /*iot*/
    name     text,             /*IoT Platform*/ 
    password text, 
);

/*context*/
DROP TABLE IF EXISTS scyna.context;
CREATE TABLE scyna.context
(
    domain      text,  /*iot*/
    code        text,  /*iot.control*/
    name        text,  /*IoT Control*/
    secret      text,
    PRIMARY KEY(domain, code)
);
CREATE INDEX ON scyna.context (code);

/*client*/
DROP TABLE IF EXISTS scyna.client;
CREATE TABLE scyna.client
(
    domain      text,  /*iot*/
    id          text,  /*datalake*/
    secret      text, 
    PRIMARY KEY(domain, id)
);
CREATE INDEX ON scyna.client (id);

/*application*/
DROP TABLE IF EXISTS scyna.application;
CREATE TABLE scyna.application
(
    domain      text, /* iot */
    code        text, /* fota */
    name        text, /* FOTA Manager */
    PRIMARY KEY(domain, code)
);
CREATE INDEX ON scyna.application (code);

/*endpoint: Command or Query*/
DROP TABLE IF EXISTS scyna.endpoint;
CREATE TABLE scyna.endpoint
(
    context     text, /* iot.control */
    url         text, /* iot/control/poweron */  
    name        text, /* PowerOn */
    description text,
    PRIMARY KEY (context, url)
);
CREATE INDEX ON scyna.endpoint (url);

/*event*/
DROP TABLE IF EXISTS scyna.event;
CREATE TABLE scyna.event
(
    context     text, /* iot.control */
    channel     text, /* iot.control.poweredon */
    name        text, /* PoweredOn */
    PRIMARY KEY (context, channel)
);
CREATE INDEX ON scyna.event (channel);

/*setting*/
DROP TABLE IF EXISTS scyna.setting;
CREATE TABLE scyna.setting
(
    context     text,
    key         text,
    value       text,
    PRIMARY KEY (context, key)
);

/*relations*/
DROP TABLE IF EXISTS scyna.client_use_endpoint;
CREATE TABLE scyna.client_use_endpoint
(
    client      text,   /*client id*/  
    endpoint     text,   /*endpoint path*/
    PRIMARY KEY (client, endpoint)
);

DROP TABLE IF EXISTS scyna.application_use_endpoint;
CREATE TABLE scyna.application_use_endpoint
(
    application     text,   /*application code*/  
    endpoint         text,   /*endpoint path*/
    PRIMARY KEY (application, endpoint)
);

DROP TABLE IF EXISTS scyna.authentication;
CREATE TABLE scyna.authentication
(
    id      text PRIMARY KEY,
    time    timestamp,
    expired timestamp,
    apps    list<text>,
    uid     text,
    data    blob,
);

DROP TABLE IF EXISTS scyna.app_has_auth;
CREATE TABLE scyna.app_has_auth
(
    app_code text,
    auth_id  text,
    user_id  text,
    PRIMARY KEY ( (app_code, user_id), auth_id )
);