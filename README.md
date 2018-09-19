===========


Due to that cloud provider usually does not support live data migration of a cloud redis service, this tool aims to dump data from one redis-server and restore the dumped data into another redis-server.


* **RESTORE** import dumped file to target redis-server

```sh
redis-dump-restore -mode=restore -host=127.0.0.1:6379 [-password=Auth] [-input=/path/to/file]
```

* **DUMP** export file from source redis-server

```sh
redis-dump-restore -mode=dump -host=127.0.0.1:6379 [-password=Auth] [-output=/path/to/file] [-database-count=16]
```

Options
-------

+ -mode=_Mode_

> Select dump mode, or restore mode. Options: Dump, Restore.

+ -host=_HostAndPort_

> The redis instance (host:port).

+ -password=_PASSWORD_

> specify the redis auth password, if empty then no use this parameter.

+ -input=_INPUT_

> use _INPUT_ as input file

+ -output=_OUTPUT_

> use _OUTPUT_ as output file

+ -database-count=_DATABASE-COUNT_

> Specify the redis database count

Examples
-------

* **RESTORE**

```sh
$ redis-dump-restore -mode=restore -input=./dump.json -host=127.0.0.1:6378
2018/09/17 23:22:30 Restored 9 Record(s).
```

* **DUMP**

```sh
$ redis-dump-restore -mode=dump -output=./dump.json
2018/09/17 23:45:55 Dumped 9 Record(s).
2018/09/17 23:45:55 Dumped 0 Record(s).
2018/09/17 23:45:55 Dumped 0 Record(s).
2018/09/17 23:45:55 Dumped 0 Record(s).
2018/09/17 23:45:55 Dumped 0 Record(s).
2018/09/17 23:45:55 Dumped 0 Record(s).
2018/09/17 23:45:55 Dumped 0 Record(s).
2018/09/17 23:45:55 Dumped 0 Record(s).
2018/09/17 23:45:55 Dumped 0 Record(s).
2018/09/17 23:45:55 Dumped 0 Record(s).
2018/09/17 23:45:55 Dumped 0 Record(s).
2018/09/17 23:45:55 Dumped 0 Record(s).
2018/09/17 23:45:55 Dumped 0 Record(s).
2018/09/17 23:45:55 Dumped 0 Record(s).
2018/09/17 23:45:55 Dumped 0 Record(s).
2018/09/17 23:45:55 Dumped 0 Record(s).

$ redis-dump-restore -mode=dump -host=127.0.0.1:6379 -output=./dump.json
2018/09/17 23:46:57 DB 0 dumped 9 Record(s).
2018/09/17 23:46:57 DB 1 dumped 0 Record(s).
2018/09/17 23:46:57 DB 2 dumped 0 Record(s).
2018/09/17 23:46:57 DB 3 dumped 0 Record(s).
2018/09/17 23:46:57 DB 4 dumped 0 Record(s).
2018/09/17 23:46:57 DB 5 dumped 0 Record(s).
2018/09/17 23:46:57 DB 6 dumped 0 Record(s).
2018/09/17 23:46:57 DB 7 dumped 0 Record(s).
2018/09/17 23:46:57 DB 8 dumped 0 Record(s).
2018/09/17 23:46:57 DB 9 dumped 0 Record(s).
2018/09/17 23:46:57 DB 10 dumped 0 Record(s).
2018/09/17 23:46:57 DB 11 dumped 0 Record(s).
2018/09/17 23:46:57 DB 12 dumped 0 Record(s).
2018/09/17 23:46:57 DB 13 dumped 0 Record(s).
2018/09/17 23:46:57 DB 14 dumped 0 Record(s).
2018/09/17 23:46:57 DB 15 dumped 0 Record(s).
```
