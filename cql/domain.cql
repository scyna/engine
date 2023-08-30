USE scyna2;

DROP TABLE IF EXISTS module;
CREATE TABLE module
(
    code        text PRIMARY KEY,
    secret      text,
);

DROP TABLE IF EXISTS endpoint;
CREATE TABLE  endpoint
(
    module      text,
    url         text,
    PRIMARY KEY (module, url)
);

DROP TABLE IF EXISTS setting;
CREATE TABLE  setting
(
    module     text,
    key         text,
    value       text,
    PRIMARY KEY (module, key)
);

DROP TABLE IF EXISTS session;
CREATE TABLE  session
(
    module      text,
    id          bigint,
    started     timestamp,
    updated     timestamp,
    ended       timestamp,
    PRIMARY KEY (module, id) 
);
CREATE INDEX ON  session (id);