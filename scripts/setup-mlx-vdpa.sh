#!/bin/bash
set -eu
 
function usage() {
        echo "$0 PCIADDR [NUM_VFs]"
        echo "  PCIADDR: The PCI address of the PF, e.g: 0000:40:00.0"
        echo "  NUM_VFS (defaul = 4): Number of VFs to configure"
        exit 1
}

function error() {
        echo $@
        exit 1
}

function get_pci_addr() {
        pf=$1
        vf=$2
        if [ -z $vf ]; then
                echo $(basename $(readlink /sys/class/net/${pf}/device))
        else
                echo $(basename $(readlink /sys/class/net/${pf}/device/virtfn${vf}))
        fi
}

# $1 is the vdpa device
# $2 is the desired vdpa driver
function set_vdpa_driver() {
        local dev=$1
        local driver=$2
            if [ -d "/sys/bus/vdpa/devices/${dev}/driver" ]; then
                    local curr_driver=$(basename $(readlink /sys/bus/vdpa/devices/${dev}/driver))
                    if [[ "${curr_driver}" != "${driver}" ]]; then
                        echo "${dev}" > /sys/bus/vdpa/drivers/${curr_driver}/unbind
                else 
                        return 0
                fi
        fi
        echo "${dev}" > /sys/bus/vdpa/drivers/${driver}/bind
}

if [ "$EUID" -ne 0 ]
  then echo "Please run as root"
  exit
fi

[ "$#" -lt 1 ] && usage

PCI_ADDR=$1
NUM_VFS=${2:-4}
 
echo "Creating VFs" 
PF=$(ls -x /sys/bus/pci/devices/${PCI_ADDR}/net/)
echo 0 > /sys/class/net/${PF}/device/sriov_numvfs
PF=$(ls -x /sys/bus/pci/devices/${PCI_ADDR}/net/)
echo 4 > /sys/class/net/${PF}/device/sriov_numvfs
PF=$(ls -x /sys/bus/pci/devices/${PCI_ADDR}/net/)

echo "Unbinding mlx5_core"
num_vfs=$(cat /sys/class/net/${PF}/device/sriov_numvfs)
for i in $(seq 0 $(($num_vfs -1))); do
        echo "Unbinding VF ${i}"
        pci_addr=$(get_pci_addr ${PF} $i)
        echo $pci_addr >  /sys/bus/pci/drivers/mlx5_core/unbind
done

echo "Modprobbing vdpa drivers"
modprobe vdpa || true
modprobe vhost-vdpa || true
modprobe virtio-vdpa || true
modprobe mlx5_vdpa || true

echo "Binding mlx5_core"
for i in $(seq 0 $(($num_vfs -1))); do
        echo "Binding VF ${i}"
        pci_addr=$(get_pci_addr ${PF} $i)
        echo $pci_addr >  /sys/bus/pci/drivers/mlx5_core/bind
        echo "Waiting for vf ${i} dev to be available "
        sleep 3
done

# Try to modprobe vhost_vdpa just in case
modprobe vhost_vdpa || true

echo ""
echo "Binding half of the devices to vhost-vdpa and other half to virtio-vpda"
DRIVERS=("vhost_vdpa" "virtio_vdpa")
devices=( $(ls -x /sys/bus/vdpa/devices) )
for i in ${!devices[@]}; do
    dev=${devices[$i]}
    driver=${DRIVERS[$(($i%2))]}
    echo "Binding device ${dev} to driver ${driver}"
    set_vdpa_driver $dev $driver
done

