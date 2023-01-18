/*event store template*/
DROP TABLE IF EXISTS module_name.event_store;
CREATE TABLE module_name.event_store
(
    event_id        bigint PRIMARY KEY,
    entity_id       bigint,
    channel         text,
    data            blob,
    time            timestamp
);

CREATE TABLE module_name.outbox
(
    event_id        bigint PRIMARY KEY,
    trace_id       bigint,
);