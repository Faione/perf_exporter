FROM docker.io/library/alpine

COPY bin/perf_exporter /app/perf_exporter

ENTRYPOINT ["/app/perf_exporter"]