DROP KEYSPACE IF EXISTS scyna;
DROP KEYSPACE IF EXISTS scyna2;
CREATE KEYSPACE scyna2 WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };