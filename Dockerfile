FROM iron/go
MAINTAINER Juraj Bubniak <juraj.bubniak@gmail.com>

COPY pgbouncer_exporter /bin

HEALTHCHECK CMD ["pgbouncer_exporter", "health"]

ENTRYPOINT ["pgbouncer_exporter"]
CMD ["server"]
