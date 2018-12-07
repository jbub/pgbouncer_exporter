# Pgbouncer exporter 
[![Build Status](https://cloud.drone.io/api/badges/jbub/pgbouncer_exporter/status.svg)][drone]
[![Docker Pulls](https://img.shields.io/docker/pulls/jbub/pgbouncer_exporter.svg?maxAge=604800)][hub]
[![Go Report Card](https://goreportcard.com/badge/github.com/jbub/pgbouncer_exporter)][goreportcard]

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

[drone]: https://cloud.drone.io/jbub/pgbouncer_exporter
[hub]: https://hub.docker.com/r/jbub/pgbouncer_exporter
[goreportcard]: https://goreportcard.com/report/github.com/jbub/pgbouncer_exporter