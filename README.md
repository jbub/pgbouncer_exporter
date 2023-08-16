# Pgbouncer exporter 
[![Build Status](https://github.com/jbub/pgbouncer_exporter/actions/workflows/go/badge.svg)][drone]
[![Docker Pulls](https://img.shields.io/docker/pulls/jbub/pgbouncer_exporter.svg?maxAge=604800)][hub]
[![Go Report Card](https://goreportcard.com/badge/github.com/jbub/pgbouncer_exporter)][goreportcard]

Prometheus exporter for Pgbouncer metrics. The minimum supported version of Pgbouncer is 1.18.

## Docker

Metrics are by default exposed on http server running on port `9127` under the `/metrics` path.

```bash
docker run \ 
  --detach \ 
  --env "DATABASE_URL=postgres://user:password@pgbouncer:6432/pgbouncer?sslmode=disable" \
  --publish "9127:9127" \
  --name "pgbouncer_exporter" \
  jbub/pgbouncer_exporter
```

In order to build the binary for the development docker compose setup you can use this command:

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
```

## Collectors

All of the collectors are enabled by default, you can control that using environment variables by settings
it to `true` or `false`.

| Name          | Description                             | Env var          | Default |
|---------------|-----------------------------------------|------------------|---------|
| stats         | Per database requests stats.            | EXPORT_STATS     | Enabled |
| pools         | Per (database, user) connection stats.  | EXPORT_POOLS     | Enabled |
| databases     | List of configured databases.           | EXPORT_DATABASES | Enabled |
| lists         | List of internal pgbouncer information. | EXPORT_LISTS     | Enabled |

## Default constant prometheus labels

In order to provide default prometheus constant labels you can use the `DEFAULT_LABELS` enviroment variable.
Labels can be set in this format `instance=pg1 env=dev`. Provided labels will be added to all the metrics. 

[hub]: https://hub.docker.com/r/jbub/pgbouncer_exporter
[goreportcard]: https://goreportcard.com/report/github.com/jbub/pgbouncer_exporter