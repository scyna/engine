USE scyna2;

DROP TABLE IF EXISTS todo;
CREATE TABLE todo
(
    bucket        bigint,
    task_id       bigint,
    PRIMARY KEY (bucket, task_id)
);

DROP TABLE IF EXISTS doing;
CREATE TABLE doing
(
    bucket        bigint,
    task_id       bigint, 
    PRIMARY KEY (bucket, task_id)
);

DROP TABLE IF EXISTS task;
CREATE TABLE task
(
    id           bigint PRIMARY KEY,
    topic        text,
    data         blob,
    start        timestamp,
    next         timestamp,
    interval     bigint,
    loop_index   bigint,
    loop_count   bigint,
    done         boolean,
);

DROP TABLE IF EXISTS module_has_task;
CREATE TABLE module_has_task
(
    module       text,
    task_id      bigint,
    PRIMARY KEY (module, task_id)
);