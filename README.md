# mackerel-plugin-pg-table-count [![Build Status](https://travis-ci.org/y-kuno/mackerel-plugin-pg-table-count.svg?branch=master)](https://travis-ci.org/y-kuno/mackerel-plugin-pg-table-count)

PostgreSQL Table Count plugin for mackerel.io agent. This repository releases an artifact to Github Releases, which satisfy the format for mkr plugin installer.

## Install

```shell
mkr plugin install y-kuno/mackerel-plugin-pg-table-count
```

## Synopsis

```shell
mackerel-plugin-postgres-table [-host=<host>] [-port=<port>] [-user=<user>] [-password=<password>] [-database=<databasename>] [-tabel=<tabel>] [-column=<column>] [-option=<option>] [-sslmode=<sslmode>] [-metric-key-prefix=<prefix>]
```

## Example of mackerel-agent.conf

```
[plugin.metrics.pg-table-count]
command = "/path/to/mackerel-plugin-postgres-table -user=postgres -database=databasename -table=tablename"
```