FROM golang:1.18.6

WORKDIR /go/src/github.com/opencurve/curve-csi
COPY . .
ARG VERSION
RUN VERSION=${VERSION} make build

FROM debian:10

USER root:root

WORKDIR /root

LABEL name="curve-csi" \
      license="Apache Version 2.0" \
      maintainers="The Netease Kubernetes Authors" \
      description="Curve CSI Plugin" \
      architecture="amd64" \
      distribution-scope="public" \
      summary="Curve CSI Plugin" \
      help="none"

RUN echo "deb http://mirrors.aliyun.com/debian/ buster main non-free contrib" > /etc/apt/sources.list && \
    echo "deb http://mirrors.aliyun.com/debian/ buster-updates main non-free contrib" >> /etc/apt/sources.list && \
    echo "deb http://mirrors.aliyun.com/debian/ buster-backports main non-free contrib" >> /etc/apt/sources.list && \
    echo "deb http://mirrors.aliyun.com/debian-security buster/updates main" >> /etc/apt/sources.list && \
    echo "deb-src http://mirrors.aliyun.com/debian/ buster main non-free contrib" >> /etc/apt/sources.list && \
    echo "deb-src http://mirrors.aliyun.com/debian/ buster-updates main non-free contrib" >> /etc/apt/sources.list && \
    echo "deb-src http://mirrors.aliyun.com/debian/ buster-backports main non-free contrib" >> /etc/apt/sources.list && \
    echo "deb-src http://mirrors.aliyun.com/debian-security buster/updates main" >> /etc/apt/sources.list

RUN apt-get update -y && \
    apt-get install -y coreutils dnsutils iputils-ping iproute2 telnet curl vim less wget graphviz unzip tcpdump gdb daemon procps python && \
    apt-get install -y ca-certificates e2fsprogs mount xfsprogs udev && \
    apt-get clean

COPY ./build/curve-csi/entrypoint.sh /root/

# Get and install curve, curve-nbd and nebd
RUN wget https://github.com/caoxianfei1/curve-csi/releases/download/v3.0.0/curve-sdk_1.2.6+72cb7582_amd64.deb && \
    wget https://github.com/caoxianfei1/curve-csi/releases/download/v3.0.0/k8s-curve-nbd_1.2.6+72cb7582_amd64.deb && \
    wget https://github.com/caoxianfei1/curve-csi/releases/download/v3.0.0/k8s-nebd_1.2.6+72cb7582_amd64.deb

# Install pkg and dependencies
RUN apt install -y libjemalloc2 && apt install libunwind8

RUN dpkg -i k8s-nebd_1.2.6+72cb7582_amd64.deb || true && apt install -fy && \
    dpkg -i k8s-curve-nbd_1.2.6+72cb7582_amd64.deb || true && apt install -fy && \
    dpkg -i curve-sdk_1.2.6+72cb7582_amd64.deb

# TODO: should created by nebd-damon but not create here
RUN mkdir -p /data/log/nebd/client && mkdir -p /data/log/nebd/server && \
    touch /data/log/nebd/nebd-server-daemon.log

RUN wget https://github.com/caoxianfei1/curve-csi/releases/download/v3.0.0/client.conf -O /etc/curve/client.conf

COPY --from=0 /go/src/github.com/opencurve/curve-csi/_output/curve-csi /bin/curve-csi

RUN chmod +x entrypoint.sh && chmod +x /bin/curve-csi
    
ENTRYPOINT ["./entrypoint.sh"]
