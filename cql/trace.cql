USE scyna2;

DROP TABLE IF EXISTS  trace;
CREATE TABLE  trace
(
    path        text,
    day         int,
    id          bigint,
    time        timestamp,
    type        smallint,
    session_id  bigint,
    duration    bigint,
    parent_id   bigint,
    source      text,
    status      int,
    PRIMARY KEY ((path, day), id)
);
CREATE INDEX ON  trace(id);

DROP TABLE IF EXISTS  client_has_trace;
CREATE TABLE  client_has_trace
(
    client_id     text,
    trace_id     bigint,
    day         int,
    PRIMARY KEY ((client_id,day), trace_id)
);

DROP TABLE IF EXISTS  span;
CREATE TABLE  span
(
    parent_id     bigint,
    child_id     bigint,
    PRIMARY KEY (parent_id, child_id)
);

DROP TABLE IF EXISTS  tag;
CREATE TABLE  tag
(
    trace_id     bigint,
    key          text,
    value        text,
    PRIMARY KEY (trace_id, key)
);

DROP TABLE IF EXISTS  log;
CREATE TABLE  log
(
    trace_id     bigint,
    seq          bigint,
    time         timestamp,
    message      text,
    level       smallint,
    PRIMARY KEY (trace_id, seq)
);
