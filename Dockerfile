FROM alpine:3.14
LABEL maintainer="Juraj Bubniak <juraj.bubniak@gmail.com>"

RUN apk --no-cache add tzdata ca-certificates

COPY pgbouncer_exporter /bin

HEALTHCHECK CMD ["pgbouncer_exporter", "health"]

ENTRYPOINT ["pgbouncer_exporter"]
CMD ["server"]
