DROP TABLE IF EXISTS scyna.trace;
CREATE TABLE scyna.trace
(
    path        text,
    day         int,
    id          bigint,
    time        timestamp,
    session_id  bigint,
    source_id   bigint,
    PRIMARY KEY ((path, day), id)
);

DROP TABLE IF EXISTS scyna.span;
CREATE TABLE scyna.span
(
    source_id     bigint,
    trace_id     bigint,
    PRIMARY KEY (source_id, trace_id)
);

DROP TABLE IF EXISTS scyna.tag;
CREATE TABLE scyna.tag
(
    trace_id     bigint,
    tag          text,
    data         blob,
    PRIMARY KEY (trace_id, tag)
);

DROP TABLE IF EXISTS scyna.log;
CREATE TABLE scyna.log
(
    trace_id     bigint,
    seq          bigint,
    time         timestamp,
    message      text,
    PRIMARY KEY (trace_id, seq)
);
