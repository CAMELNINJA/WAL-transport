
SHOW wal_level;

SELECT pg_create_logical_replication_slot('replication_slot', 'test_decoding');

SELECT slot_name, plugin, slot_type, database, active, restart_lsn, confirmed_flush_lsn FROM pg_replication_slots;


CREATE TABLE test (
  data json
);

CREATE PUBLICATION pub FOR ALL TABLES;
insert into test (data)VALUES ('{"test":"test"}')
SELECT * FROM pg_publication_tables WHERE pubname='pub';

create table t (id int, name text);
INSERT INTO t(id, name) SELECT g.id, k.name FROM generate_series(1, 10) as g(id), substr(md5(random()::text), 0, 25) as k(name);

SELECT * FROM pg_logical_slot_get_changes('replication_slot', NULL, NULL);

SELECT pg_drop_replication_slot('replication_slot');