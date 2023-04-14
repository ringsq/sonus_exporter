ARG ARCH="amd64"
ARG OS="linux"
FROM quay.io/prometheus/busybox-${OS}-${ARCH}:latest
LABEL maintainer="The RingSquared Authors "

ARG ARCH="amd64"
ARG OS="linux"
COPY .build/${OS}-${ARCH}/sonus_exporter  /bin/sonus_exporter
COPY sonus.yml       /etc/sonus_exporter/sonus.yml

EXPOSE      9700
ENTRYPOINT  [ "/bin/sonus_exporter" ]
CMD         [ "--config.file=/etc/sonus_exporter/sonus.yml" ]
