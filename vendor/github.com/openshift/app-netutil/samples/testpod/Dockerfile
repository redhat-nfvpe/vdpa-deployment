FROM centos/tools

ADD . /usr/src/app-netutil

WORKDIR /usr/src/app-netutil

ENV INSTALL_PKGS "golang"
RUN rpm --import https://mirror.go-repo.io/centos/RPM-GPG-KEY-GO-REPO && \
    curl -s https://mirror.go-repo.io/centos/go-repo.repo | tee /etc/yum.repos.d/go-repo.repo && \
    yum install -y $INSTALL_PKGS && \
    rpm -V $INSTALL_PKGS && \
    yum clean all && \
    make clean && \
    make

RUN cp /usr/src/app-netutil/bin/go_app /usr/bin/app

WORKDIR /

LABEL io.k8s.display-name="Sample application using netutil"

CMD ["/usr/bin/app", "--alsologtostderr"]
