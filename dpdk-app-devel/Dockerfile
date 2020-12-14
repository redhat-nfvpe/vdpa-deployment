FROM fedora:29
MAINTAINER Adrián Moreno <amorenoz@redhat.com>

RUN yum groupinstall -y "Development Tools"
RUN yum install -y wget numactl-devel git meson ninja-build iputils ethtool iproute

ARG repo=https://github.com/DPDK/dpdk.git
ARG version=master

ENV REPO $repo
ENV VER $version

WORKDIR /usr/src
RUN git clone $REPO dpdk

ENV DPDK_DIR=/usr/src/dpdk
WORKDIR ${DPDK_DIR}
RUN git checkout $VERSION

# Build DPDK

ENV RTE_TARGET=x86_64-native-linuxapp-gcc
ENV RTE_SDK=${DPDK_DIR}

RUN meson build; cd build; ninja

RUN cp build/app/dpdk-testpmd  /usr/bin

COPY testpmd-wrapper.sh /root
RUN chmod +x /root/testpmd-wrapper.sh

WORKDIR /root

