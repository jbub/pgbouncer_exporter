FROM alpine:3.15
LABEL maintainer="Juraj Bubniak <juraj.bubniak@gmail.com>"

RUN addgroup -S pgbouncer_exporter \
    && adduser -D -S -s /sbin/nologin -G pgbouncer_exporter pgbouncer_exporter

RUN apk --no-cache add tzdata ca-certificates

COPY pgbouncer_exporter /bin

USER pgbouncer_exporter

HEALTHCHECK CMD ["pgbouncer_exporter", "health"]

ENTRYPOINT ["pgbouncer_exporter"]
CMD ["server"]
