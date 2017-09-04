# Pgbouncer exporter [![Build Status](https://travis-ci.org/jbub/pgbouncer_exporter.svg)][travis]

[![Docker Pulls](https://img.shields.io/docker/jbub/pgbouncer_exporter.svg?maxAge=604800)][hub]
[![Go Report Card](https://goreportcard.com/badge/github.com/jbub/pgbouncer_exporter)][goreportcard]
[![Coverage Status](https://coveralls.io/repos/github/jbub/pgbouncer_exporter/badge.svg?branch=master)][coveralls]

Prometheus exporter for Pgbouncer metrics.

## Collectors

All of the collectors are enabled by default, you can control that using environment variables by settings
it to `true` or `false`.

| Name          | Description                             | Env var          | Default |
|---------------|-----------------------------------------|------------------|---------|
| stats         | Per database requests stats.            | EXPORT_STATS     | Enabled |
| pools         | Per (database, user) connection stats.  | EXPORT_POOLS     | Enabled |
| databases     | List of configured databases.           | EXPORT_DATABASES | Enabled |
| lists         | List of internal pgbouncer information. | EXPORT_LISTS     | Enabled |

[travis]: https://travis-ci.org/jbub/pgbouncer_exporter
[hub]: https://hub.docker.com/r/jbub/pgbouncer_exporter
[goreportcard]: https://goreportcard.com/report/github.com/jbub/pgbouncer_exporter
[coveralls]: https://coveralls.io/github/jbub/pgbouncer_exporter?branch=master