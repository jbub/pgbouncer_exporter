## 0.20.0

* Build with Go 1.25.
* Use alpine:3.22 as a base Docker image.
* Update dependencies.

## 0.19.0

* Add support for PgBouncer 1.24.
* Use alpine:3.21 as a base Docker image.
* Add multi-stage Docker build.

## 0.18.0

* Build with Go 1.23.
* Use alpine:3.20 as a base Docker image.
* Add server_lifetime metric.
* Use PgBouncer 1.23 in docker compose.
* Update dependencies.

## 0.17.0

* Build with Go 1.22.
* Use alpine:3.19 as a base Docker image.
* Add pool_size and max_connections metrics.
* Update dependencies.
* Update docker-compose.yml versions.

## 0.16.0

* Build with Go 1.21.
* Use alpine:3.18 as a base Docker image.
* Update github.com/prometheus/client_golang to v1.17.0.

## 0.15.0

* Build with Go 1.20.
* Use alpine:3.17 as a base Docker image.
* Add support for PgBouncer 1.18 and make it minimum supported version.
* Use Github actions instead of Drone.
* Update dependencies.

## 0.14.0

* Build with Go 1.19.
* Use alpine:3.15 as a base Docker image.
* Update dependencies.
* Update docker-compose.yml.

## 0.13.0

* Update Dockerfile to not run as nonroot.

## 0.12.0

* Build with Go 1.17.
* Use alpine:3.14 as a base Docker image.

## 0.11.0

* Add support for PgBouncer 1.16.
* Update to github.com/prometheus/common v0.30.0.

## 0.10.0

* Drop sqlx, use stdlib database.
* Add Makefile.

## 0.9.2

* Fix docker image templates.

## 0.9.1

* Use multiarch docker build to support both amd64 and arm64 platforms.

## 0.9.0

* Build with Go 1.16.
* Build also arm64 goarch.
* Bump packages.

## 0.8.0

* Use alpine:3.12 as a base Docker image.

## 0.7.0

* Add support for default constant prometheus labels.
* Bump github.com/prometheus/client_golang to v1.8.0.

## 0.6.0

* Refactor exporter to use NewConstMetric.
* Build with Go 1.15.
* Bump dependencies.

## 0.5.0

* Build with Go 1.13.
* Use sqlx.Open instead of sqlx.Connect to skip calling Ping.
* Use custom query in store.Check.
* Check store on startup.
* Add docker compose for testing.
* Update to github.com/urfave/cli/v2.
* Bump github.com/prometheus/client_golang to v1.3.0. 
* Bump github.com/lib/pq to v1.3.0.
* Update goreleaser config to support latest version.

## 0.4.0

* Build with Go 1.12.
* Pin Go modules to version tags.
* Move code to internal package.
* Switch ci from travis to drone.

## 0.3.1

* Fix build version passing in .goreleaser.yml.

## 0.3.0

* Export more metrics from stats and pools. 
* Build with Go 1.11.2.
* Add Go modules support.
* Drop dep support.

## 0.2.2

* Update vendored libs, prune tests and unused pkgs.
* Build with Go 1.10.3.
* Add golangci.yml.

## 0.2.1

* Build with Go 1.9.4.

## 0.2.0

* Add support for PgBouncer 1.8.

## 0.1.5

* Build with Go 1.9.2.
* Add Docker setup to Goreleaser config. 

## 0.1.4

* Add healthcheck.

## 0.1.3

* Refactor http server to improve testability.

## 0.1.2

* Fill missing Active field in sql store GetPools method.

## 0.1.1

* Make database column ForceUser nullable..

## 0.1.0

* Initial release.
