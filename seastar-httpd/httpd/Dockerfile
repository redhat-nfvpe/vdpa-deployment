# To build, run below cmd in the root dir of
# github.com/redhat-nfvpe/vdpa-deployment repo:
#  make httpd-image
#  -- OR --
#  make all
#  -- OR --
#  docker build --rm -t seastar-httpd -f ./seastar-httpd/httpd/Dockerfile .
#

# -------- Builder stage.
FROM fedora:29
MAINTAINER Billy McFall <bmcfall@redhat.com>

#
# Install required packages
#
RUN dnf install -y kernel-devel && \
    dnf groupinstall -y 'Development Tools'  && \
    dnf groupinstall -y 'C Development Tools and Libraries'

#
# Build httpd
#
WORKDIR /root/src/
RUN git clone https://github.com/scylladb/seastar.git

# Get the patched version of Seastar
WORKDIR /root/src/seastar
RUN git submodule update --init --recursive
RUN git remote add max https://gitlab.com/mcoquelin/seastar.git
RUN git fetch max

# Get the patched version of DPDK, sub-module to Seastar
WORKDIR /root/src/seastar/dpdk
RUN git remote add max https://gitlab.com/mcoquelin/dpdk-next-virtio.git
RUN git fetch max

# Change branches of Seastar
WORKDIR /root/src/seastar
RUN git checkout max/kubecon_poc_timer_workaround
# Lock down code to commit on 11/15/2019
RUN git checkout 7a454cc105dc24a97fa06b10124cf3850d77d4e9

# Change branches of DPDK
WORKDIR /root/src/seastar/dpdk
RUN git checkout max/seastar_kubecon_poc
# Lock down code to commit on 11/13/2019
RUN git checkout 0c8ef461ba4525bb2b6a5636cb79a5a7d10df75c

# Build
WORKDIR /root/src/seastar
RUN ./install-dependencies.sh
RUN ./configure.py --mode=dev --enable-dpdk
RUN ninja -C build/dev
RUN cp build/dev/apps/httpd/httpd /usr/bin/httpd
RUN cp build/dev/apps/seawreck/seawreck /usr/bin/seawreck

# -------- Import stage.
# BEGIN - Docker 17.05 or higher
#FROM fedora:29
#COPY --from=0 /usr/bin/httpd /usr/bin/httpd
#COPY --from=0 /usr/bin/seawreck /usr/bin/seawreck
#COPY --from=0 /lib64/libboost_atomic.so.1.66.0 /lib64/libboost_atomic.so.1.66.0
#COPY --from=0 /lib64/libboost_chrono.so.1.66.0 /lib64/libboost_chrono.so.1.66.0
#COPY --from=0 /lib64/libboost_date_time.so.1.66.0 /lib64/libboost_date_time.so.1.66.0
#COPY --from=0 /lib64/libboost_program_options.so.1.66.0 /lib64/libboost_program_options.so.1.66.0
#COPY --from=0 /lib64/libboost_system.so.1.66.0 /lib64/libboost_system.so.1.66.0
#COPY --from=0 /lib64/libboost_thread.so.1.66.0 /lib64/libboost_thread.so.1.66.0
#COPY --from=0 /lib64/libcares.so.2 /lib64/libcares.so.2
#COPY --from=0 /lib64/libcryptopp.so.7 /lib64/libcryptopp.so.7
#COPY --from=0 /lib64/libfmt.so.5 /lib64/libfmt.so.5
#COPY --from=0 /lib64/libhwloc.so.5 /lib64/libhwloc.so.5
#COPY --from=0 /lib64/libltdl.so.7 /lib64/libltdl.so.7
#COPY --from=0 /lib64/libnuma.so.1 /lib64/libnuma.so.1
#COPY --from=0 /lib64/libprotobuf.so.15 /lib64/libprotobuf.so.15
#COPY --from=0 /lib64/libyaml-cpp.so.0.6 /lib64/libyaml-cpp.so.0.6
#COPY --from=0 /lib64/ /tmp/lib64/
# END - Docker 17.05 or higher

RUN dnf install -y jq binutils && dnf clean all

LABEL io.k8s.display-name="vDPA Seastar httpd"

ADD ./seastar-httpd/httpd/entrypoint.sh /

ENTRYPOINT ["/entrypoint.sh"]
