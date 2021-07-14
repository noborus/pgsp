# pgsp - PostgreSQL Stat Progress CLI Monitor

A CUI tool that monitors PostgreSQL's pg_stat_progress*.

Supported progress reports are ANALYZE, CLUSTER, CREATE INDEX, VACUUM, and BASE_BACKUP.
See [Progress Reporting](https://www.postgresql.org/docs/current/progress-reporting.html) for more information.

![pgsp.png](https://raw.githubusercontent.com/noborus/pgsp/master/docs/pgsp.png)


## Requires

go 1.16 or later

## Install

### Download binary

[releases page](https://github.com/noborus/pgsp/releases/tag/v0.0.1).

### Go install

```console
go install github.com/noborus/pgsp/cmd/pgsp@latest
```

## Usage

Shows a progress bar if pg_stat_progress * is updated while waiting while running.

```console
$ pgsp --dsn 'host=/var/run/postgresql'
Using config file: /home/noborus/.pgsp.yaml
quit: q, ctrl+c, esc
pg_stat_progress_basebackup
 pid                  | 402006
 phase                | streaming database files
 backup_total         | 10976660480
 backup_streamed      | 6093522944
 tablespaces_total    | 1
 tablespaces_streamed | 0

█████████████████████████░░░░░░░░░░░░░░░░░░  56%
```

```console
Monitors PostgreSQL's pg_stat_progress_*.

Usage:
  pgsp [flags]

Flags:
  -a, --AfterCompletion int   Time to display after completion(Seconds) (default 10)
  -i, --Interval float        Update interval(Seconds) (default 0.5)
      --config string         config file (default is $HOME/.pgsp.yaml)
      --dsn string            PostgreSQL data source name
  -f, --fullscreen            Display in Full Screen
  -h, --help                  help for pgsp
  -t, --toggle                Help message for toggle
  -v, --version               display version information
```