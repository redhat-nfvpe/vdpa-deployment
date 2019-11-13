#!/bin/bash

# The first two commands update one of the 'if' checks to remove
# the check for 'master == VHOST_USER_SET_VRING_CALL'.
#
# Search for:   "			!(dev->flags & VIRTIO_DEV_VDPA_CONFIGURED) &&".
# Replace with: " 			!(dev->flags & VIRTIO_DEV_VDPA_CONFIGURED)) {".
sed -i -e 's/			!(dev->flags & VIRTIO_DEV_VDPA_CONFIGURED) &&/ 			!(dev->flags \& VIRTIO_DEV_VDPA_CONFIGURED)) {/g' lib/librte_vhost/vhost_user.c
#
# Search for line with: "			msg.request.master == VHOST_USER_SET_VRING_CALL) {".
# Delete the line.
sed -i -e '/			msg\.request\.master == VHOST_USER_SET_VRING_CALL) {/d' lib/librte_vhost/vhost_user.c


# Force an RARP message to be sent out.
#
# Search for line with: "	hw->started = true;".
# Append line:          "	virtio_notify_peers(dev);".
sed -i -e '/	hw->started = true;/a 	 virtio_notify_peers(dev);' drivers/net/virtio/virtio_ethdev.c
