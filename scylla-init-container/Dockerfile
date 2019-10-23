# To build, run below cmd in the root dir of
# github.com/redhat-nfvpe/vdpa-deployment repo:
#  make scylla-image
#  -- OR --
#  make all
#  -- OR --
#  docker build --rm -t scylla-init-container -f ./scylla-init-container/Dockerfile .
#

# -------- Builder stage.
FROM centos
MAINTAINER Billy McFall <bmcfall@redhat.com>

#
# Install required packages
#
RUN rpm --import https://mirror.go-repo.io/centos/RPM-GPG-KEY-GO-REPO && curl -s https://mirror.go-repo.io/centos/go-repo.repo | tee /etc/yum.repos.d/go-repo.repo
RUN yum groupinstall -y "Development Tools"
RUN yum install -y git golang make; yum clean all
# Debug Tools (if needed):
#RUN yum install -y pciutils iproute; yum clean all


#
# Build vdpadpdk-grpc-server
#
WORKDIR /root/go/src/
##RUN go get github.com/redhat-nfvpe/vdpa-deployment 2>&1 > /tmp/vdpa-deployment.log || \
##	echo "Can ignore no GO files."
ADD . /root/go/src/github.com/redhat-nfvpe/vdpa-deployment

WORKDIR /root/go/src/github.com/redhat-nfvpe/vdpa-deployment
RUN make scylla-init
RUN cp bin/scylla-init /usr/bin/scylla-init

# -------- Import stage.
# BEGIN - Docker 17.05 or higher
FROM centos
COPY --from=0 /usr/bin/scylla-init /usr/bin/scylla-init
# END - Docker 17.05 or higher

LABEL io.k8s.display-name="vDPA Scylla Init-Container"

ADD ./scylla-init-container/entrypoint.sh /

ENTRYPOINT ["/entrypoint.sh"]
