# To build:
#  docker build --rm -t uperf-app ./uperf-app
#
FROM centos:8
MAINTAINER Adrián Moreno <amorenoz@redhat.com>

RUN dnf -y install epel-release
RUN dnf -y install uperf

ADD run-uperf.sh /root
ADD iperf.xml /root
ENTRYPOINT ["/root/run-uperf.sh"]
