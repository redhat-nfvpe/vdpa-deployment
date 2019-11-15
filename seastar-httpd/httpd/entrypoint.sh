#!/bin/sh

set -e

HTTPD_SYS_BINARY_DIR="/usr/bin"
SEASTAR_DPDK_DYNAMIC_FILE="/var/run/seastar/seastar_dpdk_dynamic.conf"
SEASTAR_IP_DYNAMIC_FILE="/var/run/seastar/seastar_ip_dynamic.conf"
SEASTAR_SEAWRECK_CONFIG_FILE="/var/run/seastar/seastar_seawreck.json"

# Pod-spec can control which app to run, "httpd" or "seawreck".
# Default to "httpd" if not provided.
if [ -z $SEASTAR_APP ]; then
   echo "SEASTAR_APP not set - default to httpd"
   SEASTAR_APP=httpd
fi

if [ $SEASTAR_APP == seawreck ] ; then
   echo "SEASTAR_APP set to seawreck"
   if [ -f $SEASTAR_SEAWRECK_CONFIG_FILE ]; then
      echo "$SEASTAR_SEAWRECK_CONFIG_FILE file exists - parsing ..."
      SEASTAR_SEAWRECK_SERVER=`jq -r '.server' $SEASTAR_SEAWRECK_CONFIG_FILE`
      SEASTAR_SEAWRECK_NUM_CONNECTIONS=`jq -r '.conn' $SEASTAR_SEAWRECK_CONFIG_FILE`
      SEASTAR_SEAWRECK_DURATION_SEC=`jq -r '.d' $SEASTAR_SEAWRECK_CONFIG_FILE`
   else
      echo "$SEASTAR_SEAWRECK_CONFIG_FILE file does NOT exists - Using defaults."
      SEASTAR_SEAWRECK_SERVER="192.168.133.2:10000"
      SEASTAR_SEAWRECK_NUM_CONNECTIONS=25
      SEASTAR_SEAWRECK_DURATION_SEC=3600
   fi
   echo "Seastar-Seawreck using:"
   echo "  SEASTAR_SEAWRECK_SERVER=$SEASTAR_SEAWRECK_SERVER"
   echo "  SEASTAR_SEAWRECK_NUM_CONNECTIONS=$SEASTAR_SEAWRECK_NUM_CONNECTIONS"
   echo "  SEASTAR_SEAWRECK_DURATION_SEC=$SEASTAR_SEAWRECK_DURATION_SEC"
fi

# The associated init-container will use app-netutil to detect
# the DPDK interface passed to the container. This interface
# is passed in through a file as DPDK formatted parameters, i.e.:
#  "--vdev=net_virtio_user0,path=/var/lib/cni/usrspcni/vdpa-0,queues=1"
if [ -f $SEASTAR_DPDK_DYNAMIC_FILE ]; then
   SEASTAR_HTTPD_DPDK_DYNAMIC=`cat $SEASTAR_DPDK_DYNAMIC_FILE`
   #rm $SEASTAR_DPDK_DYNAMIC_FILE
fi

# The associated init-container will use app-netutil to detect
# the IP Address of DPDK interface passed to the container. This IP
# Address is passed in through a file.
if [ -f $SEASTAR_IP_DYNAMIC_FILE ]; then
   SEASTAR_POD_IPADDR=`cat $SEASTAR_IP_DYNAMIC_FILE`
   #rm $SEASTAR_IP_DYNAMIC_FILE
elif [ -z $SEASTAR_POD_IPADDR ]; then
	SEASTAR_POD_IPADDR="192.168.133.2"
fi


DPDK_ARGS="$SEASTAR_HTTPD_DPDK_DYNAMIC"
DPDK_ARGS="$DPDK_ARGS --no-pci"
DPDK_ARGS="$DPDK_ARGS --single-file-segments"
DPDK_ARGS="$DPDK_ARGS --iova-mode va"
DPDK_ARGS="$DPDK_ARGS -l 4"

CLI_PARAMS="--network-stack native"

if [ $SEASTAR_APP == seawreck ] ; then
   CLI_PARAMS="$CLI_PARAMS --default-log-level trace"
   CLI_PARAMS="$CLI_PARAMS --server $SEASTAR_SEAWRECK_SERVER"
   CLI_PARAMS="$CLI_PARAMS --conn $SEASTAR_SEAWRECK_NUM_CONNECTIONS"
   CLI_PARAMS="$CLI_PARAMS -d $SEASTAR_SEAWRECK_DURATION_SEC"
fi

CLI_PARAMS="$CLI_PARAMS --smp 1"
CLI_PARAMS="$CLI_PARAMS --dpdk-pmd"
CLI_PARAMS="$CLI_PARAMS --dhcp 0"
CLI_PARAMS="$CLI_PARAMS --host-ipv4-addr $SEASTAR_POD_IPADDR"
CLI_PARAMS="$CLI_PARAMS --netmask-ipv4-addr 255.255.240.0"
CLI_PARAMS="$CLI_PARAMS --collectd 0"
CLI_PARAMS="$CLI_PARAMS --hugepages /dev/hugepages"
CLI_PARAMS="$CLI_PARAMS -m 2G"
CLI_PARAMS="$CLI_PARAMS --lro off"

# Must be last so DPDK_ARGS can be appended with quotes
CLI_PARAMS="$CLI_PARAMS --argv0"

echo "SEASTAR_APP=$SEASTAR_APP"
echo "CLI_PARAMS=$CLI_PARAMS"
echo "DPDK_ARGS=\"$DPDK_ARGS\""
exec $HTTPD_SYS_BINARY_DIR/$SEASTAR_APP $CLI_PARAMS "$DPDK_ARGS"
