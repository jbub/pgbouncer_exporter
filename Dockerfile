FROM golang:1.25 AS builder
COPY . .
RUN CGO_ENABLED=0 go build -ldflags "-extldflags '-static'" -tags netgo -o /bin/pgbouncer_exporter

FROM alpine:3.22
LABEL maintainer="Juraj Bubniak <juraj.bubniak@gmail.com>"

RUN addgroup -S pgbouncer_exporter \
    && adduser -D -S -s /sbin/nologin -G pgbouncer_exporter pgbouncer_exporter

RUN apk --no-cache add tzdata ca-certificates

COPY --from=builder /bin/pgbouncer_exporter /bin

USER pgbouncer_exporter

HEALTHCHECK CMD ["pgbouncer_exporter", "health"]

ENTRYPOINT ["pgbouncer_exporter"]
CMD ["server"]
