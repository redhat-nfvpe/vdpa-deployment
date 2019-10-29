# To build, run below cmd in the root dir of
# github.com/redhat-nfvpe/vdpa-deployment repo:
#  make vdpa-cni-image
#  -- OR --
#  docker build --rm -t vdpa-cni -f ./vdpa-cni/images/Dockerfile .
#

FROM golang:alpine as builder

ADD . /usr/go/src/github.com/redhat-nfvpe/vdpa-deployment

ENV HTTP_PROXY $http_proxy
ENV HTTPS_PROXY $https_proxy
ENV GOPATH /usr/go

RUN apk add --update --virtual build-dependencies build-base linux-headers && \
    cd /usr/go/src/github.com/redhat-nfvpe/vdpa-deployment && \
    make clean && \
    make vdpa-cni

FROM alpine
COPY --from=builder /usr/go/src/github.com/redhat-nfvpe/vdpa-deployment/bin/vdpa /usr/bin/
WORKDIR /

LABEL io.k8s.display-name="vDPA CNI"

ADD ./vdpa-cni/images/entrypoint.sh /

ENTRYPOINT ["/entrypoint.sh"]
