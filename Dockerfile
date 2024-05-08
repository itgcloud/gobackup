FROM alpine:latest
ARG VERSION=latest

WORKDIR /tmp
ENV PATH="${PATH}:/opt/sqlpackage"
ARG INFLUX_CLI_VERSION=2.7.3
ARG ETCD_VER="v3.5.11"
ADD install /install

USER root

RUN <<-EOF
    set -eux
    apk add --no-cache \
    curl \
    ca-certificates \
    openssl \
    postgresql-client \
    mariadb-connector-c \
    mysql-client \
    mariadb-backup \
    redis \
    mongodb-tools \
    sqlite \
    tar \
    gzip \
    pigz \
    bzip2 \
    lzip \
    xz-dev \
    xz \
    zstd \
    libstdc++ \
    gcompat \
    icu \
    tzdata

    rm -rf /var/cache/apk/*

    wget https://aka.ms/sqlpackage-linux
    mkdir -p /opt/sqlpackage
    unzip sqlpackage-linux -d /opt/sqlpackage
    rm sqlpackage-linux
    chmod +x /opt/sqlpackage/sqlpackage

    case "$(uname -m)" in
      x86_64) arch=amd64 ;;
      aarch64) arch=arm64 ;;
      *) echo 'Unsupported architecture' && exit 1 ;;
    esac

    # Install influx
    mkdir influx && cd influx
    curl -fLO "https://dl.influxdata.com/influxdb/releases/influxdb2-client-${INFLUX_CLI_VERSION}-linux-${arch}.tar.gz" \
         -fLO "https://dl.influxdata.com/influxdb/releases/influxdb2-client-${INFLUX_CLI_VERSION}-linux-${arch}.tar.gz.asc"
    tar xzf "influxdb2-client-${INFLUX_CLI_VERSION}-linux-${arch}.tar.gz"
    cp influx /usr/local/bin/influx
    cd /tmp
    rm -rf influx
    influx version

    # Install etcdctl
    mkdir /tmp/etcd && cd /tmp/etcd
    curl -fLO https://github.com/etcd-io/etcd/releases/download/${ETCD_VER}/etcd-${ETCD_VER}-linux-${arch}.tar.gz
    tar xzf "etcd-${ETCD_VER}-linux-${arch}.tar.gz"
    cp etcd-${ETCD_VER}-linux-${arch}/etcdctl /usr/local/bin/etcdctl
    cd /tmp
    rm -rf /tmp/etcd
    etcdctl version

    /install ${VERSION} && rm /install
EOF


CMD ["/usr/local/bin/gobackup", "run"]
