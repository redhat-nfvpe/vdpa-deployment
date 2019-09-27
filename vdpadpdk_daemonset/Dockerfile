# To build:
#  docker build --rm -t nfvpe/vdpa-daemonset .
#

# -------- Builder stage.
FROM centos
MAINTAINER Billy McFall <bmcfall@redhat.com>

#
# Install required packages
#
RUN rpm --import https://mirror.go-repo.io/centos/RPM-GPG-KEY-GO-REPO && curl -s https://mirror.go-repo.io/centos/go-repo.repo | tee /etc/yum.repos.d/go-repo.repo
RUN yum groupinstall -y "Development Tools"
RUN yum install -y wget numactl-devel git golang make; yum clean all
# Debug Tools (if needed):
#RUN yum install -y pciutils iproute; yum clean all

#
# Download and Build DPDK
#
ENV DPDK_VER 19.08
ENV DPDK_DIR /usr/src/dpdk-${DPDK_VER}
WORKDIR /usr/src/
RUN wget http://fast.dpdk.org/rel/dpdk-${DPDK_VER}.tar.xz
RUN tar -xpvf dpdk-${DPDK_VER}.tar.xz

ENV RTE_TARGET=x86_64-native-linuxapp-gcc
ENV RTE_SDK=${DPDK_DIR}
WORKDIR ${DPDK_DIR}

# DPDK_VER 19.08
RUN sed -i -e 's/EAL_IGB_UIO=y/EAL_IGB_UIO=n/' config/common_linux
RUN sed -i -e 's/KNI_KMOD=y/KNI_KMOD=n/' config/common_linux
RUN sed -i -e 's/LIBRTE_KNI=y/LIBRTE_KNI=n/' config/common_linux
RUN sed -i -e 's/LIBRTE_PMD_KNI=y/LIBRTE_PMD_KNI=n/' config/common_linux
# Additional Debug if Needed
#RUN sed -i -e 's/CONFIG_RTE_LIBRTE_ETHDEV_DEBUG=n/CONFIG_RTE_LIBRTE_ETHDEV_DEBUG=y/' config/common_base

# DPDK_VER 19.02
#RUN sed -i -e 's/EAL_IGB_UIO=y/EAL_IGB_UIO=n/' config/common_linuxapp
#RUN sed -i -e 's/KNI_KMOD=y/KNI_KMOD=n/' config/common_linuxapp
#RUN sed -i -e 's/LIBRTE_KNI=y/LIBRTE_KNI=n/' config/common_linuxapp
#RUN sed -i -e 's/LIBRTE_PMD_KNI=y/LIBRTE_PMD_KNI=n/' config/common_linuxapp

RUN make install T=${RTE_TARGET} DESTDIR=${RTE_SDK}

#
# Build vdpa
#
WORKDIR ${DPDK_DIR}/examples/vdpa
COPY ./main.c ./main.c
RUN make
RUN cp build/vdpa /usr/bin/vdpa


# -------- Import stage.
# Docker 17.05 or higher
FROM centos
COPY --from=0 /usr/bin/vdpa /usr/bin/vdpa
COPY --from=0 /usr/lib64/libnuma.so.1 /usr/lib64/libnuma.so.1

LABEL io.k8s.display-name="vDPA Daemonset"

ADD ./entrypoint.sh /

ENTRYPOINT ["/entrypoint.sh"]

