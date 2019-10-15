#!/bin/bash

# Find the first string and replace it with the second string
sed -i -e 's/			!(dev->flags & VIRTIO_DEV_VDPA_CONFIGURED) &&/ 			!(dev->flags \& VIRTIO_DEV_VDPA_CONFIGURED)) {/g' lib/librte_vhost/vhost_user.c

# Find the string and delete it
sed -i -e '/			msg\.request\.master == VHOST_USER_SET_VRING_CALL) {/d' lib/librte_vhost/vhost_user.c
