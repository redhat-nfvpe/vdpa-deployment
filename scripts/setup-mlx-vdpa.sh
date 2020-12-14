#!/bin/bash
set -eu
 
usage() {
        echo "$0 PCIADDR [NUM_VFs]"
        echo "  PCIADDR: The PCI address of the PF, e.g: 0000:40:00.0"
        echo "  NUM_VFS (defaul = 4): Number of VFs to configure"
        exit 1
}

error() {
        echo $@
        exit 1
}

get_pci_addr() {
        pf=$1
        vf=$2
        if [ -z $vf ]; then
                echo $(basename $(readlink /sys/class/net/${pf}/device))
        else
                echo $(basename $(readlink /sys/class/net/${pf}/device/virtfn${vf}))
        fi
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

echo "Modprobbing mlx5_vdpa"
modprobe vdpa || true
modprobe vhost-vdpa || true
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
echo "Binding devices to vhost-vdpa driver"
for dev in $(ls -x /sys/bus/vdpa/devices); do
    if [ -d "/sys/bus/vdpa/devices/$dev/driver" ]; then
    	driver=$(basename $(readlink /sys/bus/vdpa/devices/$dev/driver))
    	if [[ "$driver" != "vhost_vdpa" ]]; then
        	echo $dev > /sys/bus/vdpa/drivers/$driver/unbind
	else
		continue
	fi
    fi  
    echo $dev > /sys/bus/vdpa/drivers/vhost_vdpa/bind
done

