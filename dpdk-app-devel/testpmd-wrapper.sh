#!/bin/bash
set -ex

RESOURCE_NAME=PCIDEVICE_${1}
shift
ARGS=$@

TESTPMD_ARGS=""
for resource in $(echo "${!RESOURCE_NAME}" | tr "," "\n"); do
        dev=$(echo $resource | cut -d ":" -f 2);
        TESTPMD_ARGS="$TESTPMD_ARGS --vdev=virtio_user0,path=$dev,mac=$MAC"; 
done

TESTPMD_ARGS="$TESTPMD_ARGS $ARGS"

ulimit -l unlimited

exec dpdk-testpmd $TESTPMD_ARGS
