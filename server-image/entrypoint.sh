#!/bin/sh

set -e

VDPA_SYS_BINARY_DIR="/usr/bin"
VDPA_SOCKETLIST="/var/run/vdpa/socketList.dat"

ORIG_NUM_WATCHES=`cat /proc/sys/fs/inotify/max_user_watches`
echo 500000 > /proc/sys/fs/inotify/max_user_watches
file=$VDPA_SOCKETLIST
while [ ! -f "$file" ]
do
    $VDPA_SYS_BINARY_DIR/inotifywait -qqt 300 -e close_write "$(dirname $file)"
done
echo $ORIG_NUM_WATCHES > /proc/sys/fs/inotify/max_user_watches


$VDPA_SYS_BINARY_DIR/vdpa-server -socketlist $VDPA_SOCKETLIST
