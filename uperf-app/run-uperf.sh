#!/bin/bash
set -e

nthr=${NTHREADS:-10}
psz=${PSIZE:-8k}
proto=${PROTO:-udp}

echo "Running uperf as $UPERF_MODE"
ip link
ip addr

if [[ $UPERF_MODE == "server" ]] ; then
    ip add add 192.168.1.1/30 dev net1
    exec uperf -sv
elif [[ $UPERF_MODE == "client" ]] ; then
    ip add add 192.168.1.2/30 dev net1
    export nthr
    export proto
    export psz
    export h=192.168.1.1
    echo "PROTOCOL ${proto} Threads: ${nthr} Packet Size: ${psz} Available CPUs $(nproc)"
    exec uperf -a -m /root/iperf.xml
else
   echo "Unknown uperf mode: $@. Only supported [server, client]"
   exit 1
fi
