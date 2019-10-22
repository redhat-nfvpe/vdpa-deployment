#!/bin/sh

set -e

VDPA_SYS_BINARY_DIR="/usr/bin"

# Let annotation updates complete before reading
sleep 5

SCYLLA_DPDK_DYNAMIC_FILE="/var/run/scylla/scylla_dpdk_dynamic.conf"

$VDPA_SYS_BINARY_DIR/scylla-init -filename $SCYLLA_DPDK_DYNAMIC_FILE
