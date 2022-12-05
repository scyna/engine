/*id*/
DROP TABLE IF EXISTS scyna.gen_id;
CREATE TABLE scyna.gen_id
(
    prefix smallint,
    seed   bigint,
    PRIMARY KEY (prefix, seed)
) WITH CLUSTERING ORDER BY (seed DESC);

/*serial number*/
DROP TABLE IF EXISTS scyna.gen_sn;
CREATE TABLE scyna.gen_sn
(
    key     text,     
    prefix  int,
    seed    int,
    PRIMARY KEY ((key,prefix), seed)
) WITH CLUSTERING ORDER BY (seed DESC);
