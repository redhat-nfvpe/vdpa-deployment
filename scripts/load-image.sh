#!/bin/bash
set -eux

# IDEA borrowed from:  https://github.com/kubernetes-sigs/kind/blob/050064b358ef09a00362ec79361081d05db4c214/pkg/cluster/nodeutils/util.go#L78

usage() {
        echo "$0 [IMAGE] [NODE]"
        echo "Loads a docker image into a kubernetes node. Only docker/containerd is supported"
        echo "  [IMAGE] The image to load"
        echo "  [NODE]  The node name as you would give it to ssh (user@server) "
        exit 1
}

save() {
    local image=$1
    local dest=$2
    echo "Saving ${image} into ${dest}"
    docker save -o ${dest} ${image}
}

load() {
    local image=$1
    local remote=$2
    echo "Loading ${image} into ${remote}"
    ssh ${remote} "docker load" < $image
}

IMAGE=$1
NODE=$2

temp=$(mktemp -d)
dest=$temp/image.tar

save ${IMAGE} ${dest}
load  ${dest} ${NODE}

rm -r ${temp}
