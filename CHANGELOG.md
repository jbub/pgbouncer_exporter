## 0.5.2

* Do not use draft release in goreleaser.

## 0.5.1

* Use docker skip_push auto in goreleaser.

## 0.5.0

* Build with Go 1.13.
* Add docker compose for testing.
* Update to github.com/urfave/cli/v2.
* Bump github.com/prometheus/client_golang to v1.3.0. 
* Bump github.com/lib/pq to v1.3.0.

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