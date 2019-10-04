#!/bin/sh

set -e

VDPA_SYS_BINARY_DIR="/usr/bin"
VDPA_PCILIST="/var/run/vdpa/pciList.dat"
VDPA_VHOST_SOCKETDIR="/var/run/vdpa/vhost"

CLI_PARAMS=""
CLI_PARAMS="$CLI_PARAMS -l 10-13"

ORIG_NUM_WATCHES=`cat /proc/sys/fs/inotify/max_user_watches`
echo 500000 > /proc/sys/fs/inotify/max_user_watches
file=$VDPA_PCILIST
while [ ! -f "$file" ]
do
    $VDPA_SYS_BINARY_DIR/inotifywait -qqt 300 -e close "$(dirname $file)"
done
echo $ORIG_NUM_WATCHES > /proc/sys/fs/inotify/max_user_watches

while read -r line
do
  CLI_PARAMS="$CLI_PARAMS -w $line,vdpa=1"
done < "$file"
rm $file

CLI_PARAMS="$CLI_PARAMS --"
CLI_PARAMS="$CLI_PARAMS --iface $VDPA_VHOST_SOCKETDIR/vdpa-"

# Remove any existing vhost socketfiles then recreated the directory
rm -rf $VDPA_VHOST_SOCKETDIR
mkdir -p $VDPA_VHOST_SOCKETDIR

$VDPA_SYS_BINARY_DIR/vdpa $CLI_PARAMS
