#!/bin/sh

set -e

VDPA_SYS_BINARY_DIR="/usr/bin"

# Let annotation updates complete before reading
sleep 5

SEASTAR_HTTPD_DPDK_DYNAMIC_FILE="/var/run/seastar/seastar_dpdk_dynamic.conf"

$VDPA_SYS_BINARY_DIR/httpd-init -filename $SEASTAR_HTTPD_DPDK_DYNAMIC_FILE
