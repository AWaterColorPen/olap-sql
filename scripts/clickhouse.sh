## CREATE TABLE
clickhouse-client --query "CREATE TABLE wikistat ( date Date, time DateTime, project String, subproject String, path String, hits UInt64, size UInt64 ) ENGINE = MergeTree(date, (path, time), 8192);"
clickhouse-client --query "CREATE TABLE wikistat_relate ( project String, class UInt64, source Float64 ) ENGINE = MergeTree() ORDER BY (project, class);"
clickhouse-client --query "CREATE TABLE wikistat_class ( id UInt64, name String ) ENGINE = MergeTree() ORDER BY (id);"
clickhouse-client --query "CREATE TABLE wikistat_uv ( date Date, time DateTime, sub_project String, uid String, activation_cnt UInt64, cost Float64 ) ENGINE = MergeTree() ORDER BY (sub_project, uid);"

## INSERT DATA TO wikistat TABLE
clickhouse-client --query "INSERT INTO wikistat VALUES ('2021-05-07','2021-05-07 11:45:26','city','CHN','level1',121,4098);"
clickhouse-client --query "INSERT INTO wikistat VALUES ('2021-05-06','2021-05-06 11:45:25','city','CHN','level1',139,10086);"
clickhouse-client --query "INSERT INTO wikistat VALUES ('2021-05-07','2021-05-07 12:43:56','city','CHN','level2',20,1024);"
clickhouse-client --query "INSERT INTO wikistat VALUES ('2021-05-07','2021-05-07 07:00:12','city','US','level1',19,2048);"
clickhouse-client --query "INSERT INTO wikistat VALUES ('2021-05-07','2021-05-07 21:23:48','school','university','engineering',2,156);"
clickhouse-client --query "INSERT INTO wikistat VALUES ('2021-05-06','2021-05-06 21:16:39','school','university','engineering',3,158);"
clickhouse-client --query "INSERT INTO wikistat VALUES ('2021-05-06','2021-05-06 20:32:41','school','senior','*',5,212);"
clickhouse-client --query "INSERT INTO wikistat VALUES ('2021-05-07','2021-05-07 09:28:27','music','pop','',4783,37291);"
clickhouse-client --query "INSERT INTO wikistat VALUES ('2021-05-07','2021-05-07 09:31:23','music','pop','ancient',391,2531);"
clickhouse-client --query "INSERT INTO wikistat VALUES ('2021-05-07','2021-05-07 09:33:59','music','rap','',1842,12942);"
clickhouse-client --query "INSERT INTO wikistat VALUES ('2021-05-07','2021-05-07 10:34:12','music','rock','',0,0);"

## INSERT DATA TO wikistat_relate TABLE
clickhouse-client --query "INSERT INTO wikistat_relate VALUES ('city',1,4.872000);"
clickhouse-client --query "INSERT INTO wikistat_relate VALUES ('school',1,0.187420);"
clickhouse-client --query "INSERT INTO wikistat_relate VALUES ('food',2,10.248400);"
clickhouse-client --query "INSERT INTO wikistat_relate VALUES ('person',3,1.730000),;"
clickhouse-client --query "INSERT INTO wikistat_relate VALUES ('music',4,93.200000);"
clickhouse-client --query "INSERT INTO wikistat_relate VALUES ('company',5,0.028100);"

## INSERT DATA TO wikistat_class TABLE
clickhouse-client --query "INSERT INTO wikistat_class VALUES (1,'location');"
clickhouse-client --query "INSERT INTO wikistat_class VALUES (2,'life');"
clickhouse-client --query "INSERT INTO wikistat_class VALUES (3,'culture');"
clickhouse-client --query "INSERT INTO wikistat_class VALUES (4,'entertainment');"
clickhouse-client --query "INSERT INTO wikistat_class VALUES (5,'social');"

## INSERT DATA TO wikistat_uv TABLE
clickhouse-client --query "INSERT INTO wikistat_uv VALUES ('2021-05-06','2021-05-06 10:00:00','CHN','aaa',0,0.000000);"
clickhouse-client --query "INSERT INTO wikistat_uv VALUES ('2021-05-06','2021-05-06 21:00:00','university','pl-okm',3,0.450000);"
clickhouse-client --query "INSERT INTO wikistat_uv VALUES ('2021-05-07','2021-05-07 07:00:00','US','aaa',4,0.510000);"
clickhouse-client --query "INSERT INTO wikistat_uv VALUES ('2021-05-07','2021-05-07 09:00:00','pop','12345678',1,0.080000);"
clickhouse-client --query "INSERT INTO wikistat_uv VALUES ('2021-05-07','2021-05-07 09:00:00','pop','qwerty',2,0.150000);"
clickhouse-client --query "INSERT INTO wikistat_uv VALUES ('2021-05-07','2021-05-07 09:00:00','rap','12345678',3,0.210000);"
clickhouse-client --query "INSERT INTO wikistat_uv VALUES ('2021-05-07','2021-05-07 10:00:00','rock','pl-okm',10,1.090000);"
clickhouse-client --query "INSERT INTO wikistat_uv VALUES ('2021-05-07','2021-05-07 11:00:00','CHN','12345678',2,0.340000);"
clickhouse-client --query "INSERT INTO wikistat_uv VALUES ('2021-05-07','2021-05-07 11:00:00','CHN','qwerty',1,0.220000);"
clickhouse-client --query "INSERT INTO wikistat_uv VALUES ('2021-05-07','2021-05-07 12:00:00','CHN','12345678',1,0.120000);"
clickhouse-client --query "INSERT INTO wikistat_uv VALUES ('2021-05-07','2021-05-07 21:00:00','university','12345678',6,1.120000);"
