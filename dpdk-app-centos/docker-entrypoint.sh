#!/bin/bash
set -e

echo "DPDK_SAMPLE_APP is set as: $DPDK_SAMPLE_APP"

if [[ $DPDK_SAMPLE_APP == l2fwd ]] ; then
   echo "Calling: l2fwd"
   exec "l2fwd"
elif [[ $DPDK_SAMPLE_APP == l3fwd ]] ; then
   echo "Calling: l3fwd"
   exec "l3fwd"
elif [[ $DPDK_SAMPLE_APP == testpmd ]] ; then
   echo "Calling: testpmd"
   exec "testpmd"
else
   echo "Net set so using default - Calling: $@"
   exec "$@"
fi
