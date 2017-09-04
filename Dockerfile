FROM iron/go
MAINTAINER Juraj Bubniak <juraj.bubniak@gmail.com>

COPY pgbouncer_exporter /bin

ENTRYPOINT ["pgbouncer_exporter"]
CMD ["server"]
