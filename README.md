===========


Due to that cloud provider usually does not support live data migration of a cloud redis service, this tool aims to dump data from one redis-server and restore the dumped data into another redis-server.


* **RESTORE** import dumped file to target redis-server

```sh
redis-transmission -mode=restore -host=127.0.0.1:6379 [-password=Auth] [-input=/path/to/file]
```

* **DUMP** export file from source redis-server

```sh
redis-transmission -mode=dump -host=127.0.0.1:6379 [-password=Auth] [-output=/path/to/file] [-database-count=16] [-thread-count=4]
```

* **SYNC** synchronize data from source redis-server to destination redis-server

```sh
redis-transmission -mode=sync -source=127.0.0.1:6379 -destination=127.0.0.1:6378 [-source-password=Auth] [-destination-password=Auth] [-database-count=16] [-sync-times=Count] [-thread-count=4]
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

+ -source=_NODE_

> The source redis instance (host:port).

+ -destination=_NODE_

> The destination redis instance (host:port).

+ -source-password=_PASSWORD_

> The source redis authorization password, if empty then no use this parameter.

+ -destination-password=_PASSWORD_

> The destination redis authorization password, if empty then no use this parameter.

+ -sync-times=_TIMES_

> synchronization times, default loop execution. Do not fill in this parameter if you need to execute it in a loop

+ -thread-count=_THREAD-COUNT_

> Number of concurrent executions, if emtpy then use cpu cores count.

+ -replace-restore=_[1|0]_

> If the destination-side not support restore command use replace option, please use 0 to off this feature, when off this feature, it will remove key before restore command executive, if empty then use replace option.

Examples
-------

* **RESTORE**

```sh
$ redis-transmission -mode=restore -input=./dump.json -host=127.0.0.1:6378
2018/09/17 23:22:30 Restored 9 Record(s).
```

* **DUMP**

```sh
$ redis-transmission -mode=dump -output=./dump.json
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

$ redis-transmission -mode=dump -host=127.0.0.1:6379 -output=./dump.json
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

* **SYNC**

```sh
$ redis-transmission -mode=sync -source=127.0.0.1:6379 -destination=127.0.0.1:6378
2018/09/20 16:42:48 Starting synchronizer
2018/09/20 16:42:48 Start 0 database thread
2018/09/20 16:42:48 Synchronized database(0) 1 records.
2018/09/20 16:42:48 Start 0 database thread
2018/09/20 16:42:48 Synchronized database(0) 1 records.
2018/09/20 16:42:48 Start 0 database thread
2018/09/20 16:42:48 Synchronized database(0) 1 records.
2018/09/20 16:42:48 Start 0 database thread
2018/09/20 16:42:48 Synchronized database(0) 1 records.
2018/09/20 16:42:48 Start 0 database thread
2018/09/20 16:42:48 Synchronized database(0) 1 records.
2018/09/20 16:42:48 Start 0 database thread
2018/09/20 16:42:48 Synchronized database(0) 1 records.
2018/09/20 16:42:48 Start 0 database thread
2018/09/20 16:42:48 Synchronized database(0) 1 records.
^C

$ redis-transmission -mode=sync -source=127.0.0.1:6379 -destination=127.0.0.1:6378 -database-count=1
2018/09/20 16:42:48 Starting synchronizer
2018/09/20 16:42:48 Start 0 database thread
2018/09/20 16:42:48 Synchronized database(0) 1 records.
2018/09/20 16:42:48 Start 0 database thread
2018/09/20 16:42:48 Synchronized database(0) 1 records.
2018/09/20 16:42:48 Start 0 database thread
2018/09/20 16:42:48 Synchronized database(0) 1 records.
2018/09/20 16:42:48 Start 0 database thread
2018/09/20 16:42:48 Synchronized database(0) 1 records.
^C

$ redis-transmission -mode=sync -source=127.0.0.1:6379 -destination=127.0.0.1:6378 -database-count=1 -sync-times=4
2018/09/20 16:42:48 Starting synchronizer
2018/09/20 16:42:48 Start 0 database thread
2018/09/20 16:42:48 Synchronized database(0) 1 records.
2018/09/20 16:42:48 Start 0 database thread
2018/09/20 16:42:48 Synchronized database(0) 1 records.
2018/09/20 16:42:48 Start 0 database thread
2018/09/20 16:42:48 Synchronized database(0) 1 records.
2018/09/20 16:42:48 Start 0 database thread
2018/09/20 16:42:48 Synchronized database(0) 1 records.
```

