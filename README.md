# pgsp - PostgreSQL Stat Progress Monitor

A CUI tool that monitors PostgreSQL's pg_stat_progress*.

## Requires

go 1.16 or later

## Install

```console
go install github.com/noborus/pgsp/cmd/pgsp@latest
```

## Usage

Shows a progress bar if pg_stat_progress * is updated while waiting while running.

```console
pgsp --dsn 'host=/var/run/postgresql'
quit: q, ctrl+c, esc
pg_stat_progress_basebackup
+-------+--------------------------+--------------+-----------------+-----------
|  PID  |          PHASE           | BACKUP TOTAL | BACKUP STREAMED | TABLESPACE
+-------+--------------------------+--------------+-----------------+-----------
| 94229 | streaming database files |  11684665344 |     10645579776 |           
+-------+--------------------------+--------------+-----------------+-----------

███████████████████████████████░░░░  91%
```