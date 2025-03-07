ARG ARCH="amd64"
ARG OS="linux"
FROM quay.io/prometheus/busybox-${OS}-${ARCH}:latest
LABEL maintainer="Nho Luong <luongutnho@hotmail.com>"

ARG ARCH="amd64"
ARG OS="linux"
COPY .build/${OS}-${ARCH}/postgres_exporter /bin/postgres_exporter

EXPOSE     9187
USER       nobody
ENTRYPOINT [ "/bin/postgres_exporter" ]
