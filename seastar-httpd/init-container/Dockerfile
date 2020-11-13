# To build, run below cmd in the root dir of
# github.com/redhat-nfvpe/vdpa-deployment repo:
#  make httpd-init-image
#  -- OR --
#  make all
#  -- OR --
#  docker build --rm -t httpd-init-container -f ./seastar-httpd/init-container/Dockerfile .
#

# -------- Builder stage.
FROM centos
MAINTAINER Billy McFall <bmcfall@redhat.com>

#
# Install required packages
#
RUN yum groupinstall -y "Development Tools"
RUN yum install -y git golang make; yum clean all
# Debug Tools (if needed):
#RUN yum install -y pciutils iproute; yum clean all


#
# Build Seastar-httpd Init Code
#
WORKDIR /root/go/src/
##RUN go get github.com/redhat-nfvpe/vdpa-deployment 2>&1 > /tmp/vdpa-deployment.log || \
##	echo "Can ignore no GO files."
ADD . /root/go/src/github.com/redhat-nfvpe/vdpa-deployment

WORKDIR /root/go/src/github.com/redhat-nfvpe/vdpa-deployment
RUN make httpd-init
RUN cp bin/httpd-init /usr/bin/httpd-init

# -------- Import stage.
# BEGIN - Docker 17.05 or higher
FROM centos
COPY --from=0 /usr/bin/httpd-init /usr/bin/httpd-init
# END - Docker 17.05 or higher

LABEL io.k8s.display-name="vDPA Seastart-httpd Init-Container"

ADD ./seastar-httpd/init-container/entrypoint.sh /

ENTRYPOINT ["/entrypoint.sh"]
