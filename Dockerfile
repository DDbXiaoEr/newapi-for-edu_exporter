FROM debian:bullseye-v1.4.1

COPY bin/newapi-exporter /app/newapi-exporter

EXPOSE 9090

ENTRYPOINT ["/app/newapi-exporter"]
