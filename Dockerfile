FROM alpine:3.14
LABEL maintainer="Juraj Bubniak <juraj.bubniak@gmail.com>"

RUN addgroup -S pgexporter \
    && adduser -D -S -s /sbin/nologin -G pgexporter pgexporter \
    && apk --no-cache add tzdata ca-certificates

COPY pgbouncer_exporter /bin

USER pgexporter

HEALTHCHECK CMD ["pgbouncer_exporter", "health"]

ENTRYPOINT ["pgbouncer_exporter"]
CMD ["server"]
