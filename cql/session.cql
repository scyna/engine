DROP TABLE IF EXISTS scyna.session;
CREATE TABLE scyna.session
(
    id          bigint,
    context     text,
    start       timestamp,
    last_update timestamp,
    end         timestamp,
    exit_code   int, /*normal, crash, ....*/
    PRIMARY KEY (context, id) 
);
CREATE INDEX ON scyna.session (id);

DROP TABLE IF EXISTS scyna.session_log;
CREATE TABLE scyna.session_log
(
    session_id bigint,
    day        int,
    time       timestamp,
    seq        bigint,
    message    text,
    level      smallint,
    PRIMARY KEY ((session_id, day), seq) /*session_has_log*/
);

