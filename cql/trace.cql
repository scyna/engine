USE scyna2;

DROP TABLE IF EXISTS trace;
CREATE TABLE trace
(
    path        text,
    day         int,
    time        timestamp,
    id          bigint,
    type        smallint,
    session     bigint,
    duration    bigint,
    parent      bigint,
    status      int,
    PRIMARY KEY ((path,day), time, id)
) WITH CLUSTERING ORDER BY (time DESC);
CREATE INDEX ON  trace(id);
CREATE INDEX ON  trace(parent);

DROP TABLE IF EXISTS endpoint_trace;
CREATE TABLE trace
(
    day         int,
    time        timestamp,
    id          bigint,
    request     text,
    response    text,
    status      int,
    session     bigint,
    PRIMARY KEY (day, time, id)
) WITH CLUSTERING ORDER BY (time DESC);
CREATE INDEX ON endpoint_trace(id);

DROP TABLE IF EXISTS log;
CREATE TABLE log
(
    day          int,
    time         timestamp,
    source       bigint,  /*trace_id or session_id*/
    seq          bigint,
    message      text,
    level        smallint,
    PRIMARY KEY (day, time, source, seq)
) WITH CLUSTERING ORDER BY (time DESC);
CREATE INDEX ON log(source);