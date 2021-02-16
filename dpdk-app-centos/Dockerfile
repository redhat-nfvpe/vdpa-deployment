# To build:
#  docker build --rm -t dpdk-app-centos ./dpdk-app-centos
#


# -------- Builder stage.
FROM centos:8
MAINTAINER Billy McFall <bmcfall@redhat.com>

#
# Install required packages
#
#RUN rpm --import https://mirror.go-repo.io/centos/RPM-GPG-KEY-GO-REPO && curl -s https://mirror.go-repo.io/centos/go-repo.repo | tee /etc/yum.repos.d/go-repo.repo
RUN dnf groupinstall -y "Development Tools"
RUN dnf module -y install go-toolset
RUN dnf install -y dnf-plugins-core; dnf -y config-manager --set-enabled powertools
RUN dnf install -y wget numactl-devel git golang meson ninja-build make
# Needed by latest DPDK
RUN pip3 install pyelftools

# Debug Tools (if needed):
#RUN dnf install -y pciutils iproute; yum clean all
# Uncomment to build DPDK with debug symbols
#ENV MESONOPTS="--buildtype=debug"

#
# Download and Build APP-NetUtil
#
WORKDIR /root/go/src/
RUN go get github.com/openshift/app-netutil 2>&1 > /tmp/UserspaceDockerBuild.log || echo "Can ignore no GO files."
#RUN go get github.com/openshift/app-netutil
WORKDIR /root/go/src/github.com/openshift/app-netutil
RUN make c_sample
RUN cp bin/libnetutil_api.so /lib64/libnetutil_api.so; cp bin/libnetutil_api.h /usr/include/libnetutil_api.h

#
# Download and Build DPDK
# Uncomment to select stable branch
#ENV DPDK_VER 20.11
#ENV DPDK_DIR /usr/src/dpdk-${DPDK_VER}
#WORKDIR /usr/src/
#RUN wget http://fast.dpdk.org/rel/dpdk-${DPDK_VER}.tar.xz
#RUN tar -xpvf dpdk-${DPDK_VER}.tar.xz

# Uncomment to select upstream branch
ENV DPDK_BRANCH=main
ENV DPDK_DIR /usr/src/dpdk
RUN git clone --branch ${DPDK_BRANCH} https://github.com/dpdk/dpdk.git ${DPDK_DIR}

WORKDIR ${DPDK_DIR}

#
# Substitute Testpmd
#
WORKDIR ${DPDK_DIR}/app/test-pmd
COPY ./dpdk-args.c ./dpdk-args.c
COPY ./dpdk-args.h ./dpdk-args.h
COPY ./testpmd_eal_init.txt ./testpmd_eal_init.txt
COPY ./testpmd_launch_args_parse.txt ./testpmd_launch_args_parse.txt
COPY ./testpmd_substitute.sh ./testpmd_substitute.sh
RUN ./testpmd_substitute.sh

#
# Substitute l2fwd
#
WORKDIR ${DPDK_DIR}/examples/l2fwd
COPY ./dpdk-args.c ./dpdk-args.c
COPY ./dpdk-args.h ./dpdk-args.h
COPY ./l2fwd_eal_init.txt ./l2fwd_eal_init.txt
COPY ./l2fwd_parse_args.txt ./l2fwd_parse_args.txt
COPY ./l2fwd_substitute.sh ./l2fwd_substitute.sh
RUN ./l2fwd_substitute.sh

#
# Substitute l3fwd
#
WORKDIR ${DPDK_DIR}/examples/l3fwd
COPY ./dpdk-args.c ./dpdk-args.c
COPY ./dpdk-args.h ./dpdk-args.h
COPY ./l3fwd_eal_init.txt ./l3fwd_eal_init.txt
COPY ./l3fwd_parse_args.txt ./l3fwd_parse_args.txt
COPY ./l3fwd_substitute.sh ./l3fwd_substitute.sh
RUN ./l3fwd_substitute.sh


#
# Build
#
WORKDIR ${DPDK_DIR}
RUN meson -Dexamples=l3fwd,l2fwd ${MESONOPTS} build && cd build && ninja

RUN cp build/examples/dpdk-l3fwd /usr/bin/l3fwd
RUN cp build/examples/dpdk-l2fwd /usr/bin/l2fwd
RUN cp build/app/dpdk-testpmd /usr/bin/testpmd

# -------- Import stage.
# Docker 17.05 or higher
# BEGIN
FROM centos:8
COPY --from=0 /usr/bin/l2fwd /usr/bin/l2fwd
COPY --from=0 /usr/bin/l3fwd /usr/bin/l3fwd
COPY --from=0 /usr/bin/testpmd /usr/bin/testpmd
COPY --from=0 /lib64/libnetutil_api.so /lib64/libnetutil_api.so
COPY --from=0 /usr/lib64/libnuma.so.1 /usr/lib64/libnuma.so.1
# Set default app
RUN ln -s /usr/bin/l3fwd /usr/bin/dpdk-app
# END

COPY ./docker-entrypoint.sh /
ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["dpdk-app"]
