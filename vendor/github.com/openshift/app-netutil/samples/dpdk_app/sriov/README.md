#  SR-IOV Deployment with dpdk-app-centos image
This directory contains the YAML files needed to start the SR-IOV
Device Plugin as a daemonset and launch a DPDK based docker image
which leverages the SR-IOV CNI to plug SR-IOV VFs into the pod. This
file only describes how these files are used. See the following two
repos for more details on how SR-IOV Device Plugin and CNI work and
how to build the SR-IOV Device Plugin image (SR-IOV Device Plugin is
run in a container as a daemonset) and build the SR-IOV CNI:
* https://github.com/intel/sriov-network-device-plugin
* https://github.com/intel/sriov-cni

The files to build the `dpdk-app-centos` image are located in this
same repo. See the following link for how to build the image and
details regarding what the image is doing:
* [dpdk-app-centos](../dpdk-app-centos/README.md)

# SR-IOV Setup
This test setup assumes:
* Running on baremetal.
* Kubernetes, Multus and SR-IOV CNI are installed.
* SR-IOV VFs have already been created on the PFs being used.

This test setup uses two physical NICs (PFs) with one VF from each PF
attached to the pod. It maps traffic from the PF to the each VF using
VLANs. 

## Download the sample yaml files
Use the following steps to download the sample YAML files:
```
cd $GOPATH/src
go get github.com/openshift/app-netutil
cd github.com/openshift/app-netutil/samples/dpdk_app/sriov/
```

The following sections all assume your working directory
is `$GOPATH/src/github.com/openshift/app-netutil/samples/dpdk_app/sriov/`.

## Create the Network-Attachment-Definition for each desired network
This setup assumes there are two networks, one network for each PF. The
following commands setup those networks:
```
kubectl create -f netAttach-sriov-dpdk-a.yaml
kubectl create -f netAttach-sriov-dpdk-b.yaml
```

These YAML files map a VLAN to the VF. It is currently using VLAN 100 for
"sriov-network-a" and VLAN 200 for "sriov-network-b". These values can be
changed if needed.

The following command can be used to determine the set of
Network-Attachment-Definitions currently created on the system:
```
kubectl get network-attachment-definitions
NAME          AGE
sriov-net-a   4h18m
sriov-net-b   4h18m
```

## Create ConfigMap
The following command creates the configMap. The ConfigMap provides the
filters to the SR-IOV Device-Plugin to allow it to select the set of VFs
that are available to a given Network-Attachment-Definition. This file
uses the PFs `eno1` and `eno2`. If your system is using other interfaces,
then update the file accordingly.

`NOTE:` This file will most likely need to be updated before using.

```
kubectl create -f ./configMap.yaml
```

The following command can be used to collect info on which interfaces are
in your system and manufacturer details that are used in the configMap to
select available VFs. 
```
lspci -nn | grep Ethernet
01:00.0 Ethernet controller [0200]: Intel Corporation Ethernet Controller X710 for 10GbE SFP+ [8086:1572] (rev 01)
01:00.1 Ethernet controller [0200]: Intel Corporation Ethernet Controller X710 for 10GbE SFP+ [8086:1572] (rev 01)
01:02.0 Ethernet controller [0200]: Intel Corporation Ethernet Virtual Function 700 Series [8086:154c] (rev 01)
01:02.1 Ethernet controller [0200]: Intel Corporation Ethernet Virtual Function 700 Series [8086:154c] (rev 01)
:
```

The following command can be used to determine the set of configMaps
currently created in the system:
```
kubectl get configmaps  --all-namespaces
NAMESPACE     NAME                                 DATA   AGE
kube-public   cluster-info                         2      5d23h
kube-system   coredns                              1      5d23h
kube-system   extension-apiserver-authentication   6      5d23h
kube-system   kube-flannel-cfg                     2      5d23h
kube-system   kube-proxy                           2      5d23h
kube-system   kubeadm-config                       2      5d23h
kube-system   kubelet-config-1.15                  1      5d23h
kube-system   multus-cni-config                    1      5d23h
kube-system   sriovdp-config                       1      4h24m
```

## Start the SR-IOV Device Plugin Daemonset
The following command starts the SR-IOV Device Plugin as a
daemonset container:  
```
kubectl create -f sriovdp-daemonset.yaml
```

To determine if the SR-IOV Device Plugin is running, use the
following command and find the
`kube-sriov-device-plugin-amd64-xxx` pod:
```
kubectl get pods --all-namespaces
NAMESPACE     NAME                                    READY   STATUS    RESTARTS   AGE
kube-system   coredns-5c98db65d4-78v6k                1/1     Running   16         5d23h
kube-system   coredns-5c98db65d4-r5mmj                1/1     Running   16         5d23h
kube-system   etcd-nfvsdn-22-oot                      1/1     Running   16         5d23h
kube-system   kube-apiserver-nfvsdn-22-oot            1/1     Running   16         5d23h
kube-system   kube-controller-manager-nfvsdn-22-oot   1/1     Running   16         5d23h
kube-system   kube-flannel-ds-amd64-jvnm5             1/1     Running   16         5d23h
kube-system   kube-multus-ds-amd64-lxv5v              1/1     Running   16         5d23h
kube-system   kube-proxy-6w7sn                        1/1     Running   16         5d23h
kube-system   kube-scheduler-nfvsdn-22-oot            1/1     Running   16         5d23h
kube-system   kube-sriov-device-plugin-amd64-6cj7g    1/1     Running   0          4h6m
```

Once the SR-IOV Device Plugin is started, it probes the system
looking for VFs that meet the selector’s criteria. This takes a
couple of seconds to collect. The following command can be used to
determine the number of detected VFs. (NOTE: This is the allocated
values and does not change as VFs are doled out.) See
"intel.com/intel_sriov_dpdk_a" and "intel.com/intel_sriov_dpdk_b":
```
kubectl get node nfvsdn-22-oot -o json | jq '.status.allocatable'
{
  "cpu": "64",
  "ephemeral-storage": "396858657750",
  "hugepages-1Gi": "64Gi",
  "intel.com/intel_sriov_dpdk_a": "8",
  "intel.com/intel_sriov_dpdk_b": "8",
  "memory": "64773512Ki",
  "pods": "110"
}
```

## Start the DPDK based container
Use the following command to start the DPDK based container using
SR-IOV Interfaces:
```
kubectl create -f sriov-pod-1.yaml
```

If needed, ‘exec’ into the container to customize DPDK application:
```
kubectl exec -it sriov-pod-1 -- sh
```

## Test Generator
By default, the DPDK based container ‘dpdk-app-centos’ is running the DPDK
‘l3fwd’ sample application (see https://doc.dpdk.org/guides/sample_app_ug/l3_forward.html).
This sample application does some simple routing based on a hard-coded routing
table. The following subnets are assigned to interfaces:
```
Interface 0: Route 192.18.0.0 / 24
Interface 1: Route 192.18.1.0 / 24
Interface 2: Route 192.18.2.0 / 24
Interface 3: Route 192.18.3.0 / 24
Interface 4: Route 192.18.4.0 / 24
:
```

In the test setup described in this document, VFs from the first interface, ‘eno1’,
will be assigned ‘Interface 0’ in DPDK, and thus will get route `192.18.0.0 / 24`
assigned to it. VFs from the second interface, ‘eno2’, will be assigned ‘Interface 1’
in DPDK, and thus will get route 192.18.1.0 / 24 assigned to it. At this time, this
is not configurable.

As described above, VLAN IDs are used to map packets from the PF to a given VF. This
value is configurable, but in test setup described here, VLAN 100 is used to map packets
from the first PF, ‘eno1’, to its associated VF. VLAN 200 is used to map packets from
the second PF, ‘eno2’, to its associated VF.

# SR-IOV Teardown
The following steps are used to stop the container and SR-IOV Device
Plugin:
```
kubectl delete pod sriov-pod-1
kubectl delete -f sriovdp-daemonset.yaml
kubectl delete -f configMap.yaml
kubectl delete -f netAttach-sriov-dpdk-b.yaml
kubectl delete -f netAttach-sriov-dpdk-a.yaml
```