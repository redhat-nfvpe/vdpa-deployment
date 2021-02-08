# vdpa-deployment
Example YAML files to deploy vDPA VFs in a container running in kubernetes.

## Overview
VirtIO Data Path Acceleration (vDPA) is a technology that enables pods to use
accelerated network interfaces without having to include vendor specific
drivers. This is possible because vDPA-capable NICs implement the virtIO
datapath. The vDPA Framework is in charge of translating the vendor-specific
control path (that the NIC understands) to a vendor agnostic protocol
(to be exposed to the application).

For an overview of the technology, read the
[vDPA overview blog post](https://www.redhat.com/en/blog/introduction-vdpa-kernel-framework).
More technical blog entries can also be read in the
[Virtio-networking series two](https://www.redhat.com/en/blog/virio-networking-series-advanced).

Note that, apart from the vDPA kernel framework implemented in the
linux kernel, there is another vDPA framework in DPDK. However, the DPDK
framework is out of the scope of this repository for now.

This repo combines several other repos to enable vDPA VFs to be used in
containers. The following diagram shows an overview of the end-to-end
vDPA solution in Kubernetes:


![](doc/images/vDPA_SRIOV_design-Legacy-kernel.png)

More information about this solution can be found in the [Design Document](https://docs.google.com/document/d/1DgZuksLVIVD5ZpNUNH7zPUr-8t6GKKQICDLqIwQv-FA)

As shown in the diagram, the Kubernetes vDPA solution will support both
SR-IOV CNI (for legacy SR-IOV devices) and the [Accelerated Bridge CNI](https://github.com/Mellanox/accelerated-bridge-cni/)
(for switchdev devices). Currently, this repository focuses on
using SR-IOV CNI

## Quick Start
To leverage this repo, download this repo, run `make all`:

```
   make all
```

`make all` builds the following images/binaries:

* `sriov-device-plugin` docker image: Located in the **sriov-dp**
  directory. This image takes the upstream SR-IOV Device Plugin and applies
  some local patches to enable it to work with vDPA as well. See
  [sriov-dp](#sriov-dp).
* `sriov-cni` binary and docker image: Located in the **sriov-cni** directory.
  To install the sriov-cni, the binary must be copied to the default CNI directory,
  typically `/opt/cni/bin/`. Alternatively, a DaemonSet can be deployed which will
  take care of doing that in all the nodes. See [sriov-cni](#sriov-cni).
* `dpdk-app-devel` docker image: This image contains a recent DPDK installation and
some development utilities
* `dpdk-app` docker image: This image contains a centos-based DPDK application powered
by [app-netutil](https://github.com/openshift/app-netutil) that is able to run testpmd,
l2fwd and l3fwd
* `multus` image: This image deploys a multus binary that supports the
[Device-info spec](https://github.com/k8snetworkplumbingwg/device-info-spec).

If you don't want to build all the projects from source, docker images
will be provided for convenience. See [Docker Hub section](#docker-hub)

On multi-node clusters you might need to load the built images
into the different nodes:
```
    ./scripts/load-image.sh nfvpe/sriov-device-plugin user@HOSTNAME
    ./scripts/load-image.sh nfvpe/sriov-cni user@HOSTNAME
    ./scripts/load-image.sh nfvpe/multus user@HOSTNAME
    ./scripts/load-image.sh dpdk-app-centos user@HOSTNAME

```
The following set of commands will deploy the images above.
```
    make deploy
```

Update configMap-vdpa.yaml to match local HW and then run
```
    kubectl create -f ./deployment/configMap-vdpa.yaml
```

Finally, some sample network-attachment-definitions are available in
`deployment`

To deploy a sample application, see [Sample Applications](#sample-applications)


## Sample Applications
Once the SR-IOV Device Plugin and SR-IOV CNI have been installed, the application
consuming the vDPA devices can be started. This repository will provide some
sample applications:

* **single-pod**: A simgle DPDK pod using a vDPA interface
* **vdpa-traffic-test**: A simple test that deploys two pods that send packets to each
    other (using testpmd)
* More TBD

### single-pod
The single pod application deploys a pod that runs testpmd on the vdpa device.
The testpmd arguments can be modified in `deployments/vdpa-single.yaml`

To deploy the application run:

```
    kubectl apply -f deployment/vdpa-single.yaml
```

Inspect the logs with:

```
    kubectl apply logs -f vdpa-pod
```

Delete the application by running:


```
    kubectl delete -f deployment/vdpa-single.yaml
```

### vdpa-dpdk-traffic-test
The traffic test deploys two pods using dpdk and
[app-netutil](https://github.com/openshift/app-netutil). One generates traffic
and the other receives it.

In order to select where the generator and sink runs, node selectors are used.

First, add a label to the node you want the generator to run on:

```
    kubectl label node GEN_NODENAME vdpa-test-role-gen=true
```

```
    kubectl label node SINK_NODENAME vdpa-test-role-sink=true
```


Deploy the application by running:

```
    kubectl apply -f deployment/netAttach-vdpa-vhost-mlx-1000.yaml
    kubectl apply -f deployment/netAttach-vdpa-vhost-mlx-2000.yaml
    kubectl apply -f deployment/vdpa-dpdk-traffic-test.yaml
```

Delete the application by running:

```
    kubectl delete -f deployment/vdpa-dpdk-traffic-test.yaml
    kubectl delete -f deployment/netAttach-vdpa-vhost-mlx-1000.yaml
    kubectl delete -f deployment/netAttach-vdpa-vhost-mlx-2000.yaml
```

Two env variables can be used in the deployment file
(`deployment/vdpa-dpdk-traffic-test.yaml`) to modify the behavior of the pod.

* **DPDK_SAMPLE_APP** can be used to select between testpmd, l2fwd and l3fwd.
* **TESTPMD_EXTRA_ARGS** can be used to add extra command line arguments to 
testpmd. For example "--forward-mode=flowgen". Those arguments will be appended
to the default ones: "--auto-start --tx-first --stats-period 2"


### Prerequisites
This setup assumes:
* Running on bare metal.
* Kubernetes is installed.
* Multus CNI is installed.
* vDPA VFs have already been created and bound to vhost-vdpa driver

For reference, this repo was developed and tested on:
* Fedora 32 - kernel 5.10.0+ (modified: See [HugePage Cgroup Known Issue](#hugepage-cgroup))
* GO: go1.14.9
* Docker: 19.03.11
* Kubernetes: v1.19.3


### Details

### Supported Hardwared
This repo has been tested with:
- Nvidia Mellanox ConnectX-6 Dx


### Kubernetes-vDPA setup

To deploy the Kubernetes-vDPA solution, the following steps must be taken:
* Install SR-IOV CNI
* Create Network-Attachment-Definition
* Create ConfigMap
* Start SR-IOV Device Plugin Daemonset


### sriov-cni
The changes to enable the SR-IOV CNI to also manage vDPA interfaces
are in this repository:

https://github.com/amorenoz/sriov-cni/tree/rfe/vdpa

#### sriov-cni Image

To build SR-IOV CNI in a Docker image:
```
   make sriov-cni
```

To run:
```
   kubectl create -f ./deployment/sriov-cni-daemonset.yaml
```

As with all DaemonSet YAML files, there is a version of the file for
Kubernetes versions prior to 1.16 in the `k8s-pre-1-16` subdirectory.


#### Network-Attachment-Definition
The Network-Attachment-Definition define the attributes of the network
for the interface (in this case a vDPA VF) that is being attached to
the pod.

There are three sample Network-Attachment-Definition in the `deployment`
directory. You can modify them freely to match your setup. For more
information, see the [SR-IOV CNI Configuration reference](https://github.com/k8snetworkplumbingwg/sriov-cni/blob/master/docs/configuration-reference.md)

The following commands setup those networks:
```
   kubectl create -f ./deployment/netAttach-vdpa-vhost-mlx.yaml
   kubectl create -f ./deployment/netAttach-vdpa-vhost-mlx-1000.yaml
   kubectl create -f ./deployment/netAttach-vdpa-vhost-mlx-2000.yaml
```

The following command can be used to determine the set of
Network-Attachment-Definitions currently created on the system:
```
  kubectl get network-attachment-definitions
  NAME                      AGE
  vdpa-mlx-vhost-net        24h
  vdpa-mlx-vhost-net-1000   24h
  vdpa-mlx-vhost-net-2000   24h
```

The following commands delete those networks:
```
   kubectl delete -f ./deployment/netAttach-vdpa-vhost-mlx.yaml
   kubectl delete -f ./deployment/netAttach-vdpa-vhost-mlx-1000.yaml
   kubectl delete -f ./deployment/netAttach-vdpa-vhost-mlx-2000.yaml
```


#### ConfigMap
The ConfigMap provides the filters to the SR-IOV Device-Plugin to
allow it to select the set of VFs that are available to a given 
Network-Attachment-Definition. The parameter ‘resourceName’ maps
back to one of the Network-Attachment-Definitions defined earlier.

The SR-IOV Device Plugin has been extended to support an additional
filter that is used to select the vdpa type to be used: `vdpaType`.
Supported values are:

* vhost
* virtio

The following example configMap creates two pools of vdpa devices
bound to vhost-vdpa driver:

Example:
```
cat deployment/configMap-vdpa.yaml 
apiVersion: v1
kind: ConfigMap
metadata:
  name: sriovdp-config
  namespace: kube-system
data:
  config.json: |
    {
        "resourceList": [{
                "resourceName": "vdpa_ifcvf_vhost",
                "selectors": {
                    "vendors": ["1af4"],
                    "devices": ["1041"],
                    "drivers": ["ifcvf"],
                    "vdpaType": "vhost"
                }
            },
            {
                "resourceName": "vdpa_mlx_vhost",
                "selectors": {
                    "vendors": ["15b3"],
                    "devices": ["101e"],
                    "drivers": ["mlx5_core"],
                    "vdpaType": "vhost"
                }
            }
        ]
    }

```

`NOTE:` This file will most likely need to be updated before using to
match interface on deployed hardware. To obtain the other attributes,
like vendor and devices, use the ‘lspci’ command:
```
lspci -nn | grep Ethernet
:
05:00.1 Ethernet controller [0200]: Intel Corporation Device [8086:15fe]
05:00.2 Ethernet controller [0200]: Red Hat, Inc. Virtio network device [1af4:1041] (rev 01)
05:00.3 Ethernet controller [0200]: Red Hat, Inc. Virtio network device [1af4:1041] (rev 01)
65:00.0 Ethernet controller [0200]: Mellanox Technologies MT2892 Family [ConnectX-6 Dx] [15b3:101d]
65:00.1 Ethernet controller [0200]: Mellanox Technologies MT2892 Family [ConnectX-6 Dx] [15b3:101d]
65:00.2 Ethernet controller [0200]: Mellanox Technologies ConnectX Family mlx5Gen Virtual Function [15b3:101e]
65:00.3 Ethernet controller [0200]: Mellanox Technologies ConnectX Family mlx5Gen Virtual Function [15b3:101e]
```

The following command creates the configMap:
```
   cd $GOPATH/src/github.com/redhat-nfvpe/vdpa-deployment
   kubectl create -f ./deployment/configMap-vdpa.yaml
```

The following command can be used to determine the set of
configMaps currently created in the system:
```
kubectl get configmaps  --all-namespaces
NAMESPACE     NAME                                 DATA   AGE
kube-public   cluster-info                         2      5d23h
kube-system   coredns                              1      5d23h
kube-system   extension-apiserver-authentication   6      5d23h
kube-system   multus-cni-config                    1      5d23h
kube-system   sriovdp-config                       1      4h24m
```

The following command deletes the configMap:
```
   kubectl delete -f ./deployment/configMap-vdpa.yaml
```

#### SR-IOV Device Plugin DaemonSet
The changes to enable the SR-IOV Device Plugin to also manage
vDPA interfaces are currently in this repository:

https://github.com/amorenoz/sriov-network-device-plugin/tree/vdpaInfoProvider

To build the SR-IOV Device Plugin run:

```
   make sriov-dp
```

To build from scratch:

```
   make sriov-dp SCRATCH=y
```

Deploy the SR-IOV Device Plugin by running the following command:

```
   kubectl create -f ./deployment/sriov-dp-daemonset.yaml
```

#### SR-IOV Device Plugin DaemonSet
The SR-IOV Device Plugin runs as a DaemonSet (always running as opposed
to CNI which is called and returns immediately). It is recommended that
the SR-IOV Device Plugin run in a container. So this set is to start the
container the SR-IOV Device Plugin is running in.

The following command started the SR-IOV Device Plugin DaemonSet:
```
   cd $GOPATH/src/github.com/redhat-nfvpe/vdpa-deployment
   kubectl create -f ./deployment/sriov-vdpa-daemonset.yaml
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

```
for node in $(kubectl get nodes | grep Ready | awk '{print $1}' ); do echo "Node $node:" ; kubectl get node $node -o json | jq '.status.allocatable'; done
Node virtlab711.virt.lab.eng.bos.redhat.com:
{
  "cpu": "32",
  "ephemeral-storage": "859332986687",
  "hugepages-1Gi": "10Gi",
  "hugepages-2Mi": "0",
  "intel.com/vdpa_intel_vhost: "0",
  "intel.com/vdpa_mlx_vhost": "2",
  "memory": "120946672Ki",
  "pods": "110"
}
Node virtlab712.virt.lab.eng.bos.redhat.com:
{
  "cpu": "32",
  "ephemeral-storage": "844837472087",
  "hugepages-1Gi": "10Gi",
  "hugepages-2Mi": "0",
  "intel.com/vdpa_intel_vhost: "0",
  "intel.com/vdpa_mlx_vhost": "2",
  "memory": "120950288Ki",
  "pods": "110"
}
```


## Docker Hub
All the images have been pushed to Docker Hub. TBD

## Known issues and limitations
### Hugepage Cgroup
There is an issue inrecent kernels (>=5.7.0) that affects hugetlb-cgroup reservation.
There are two ways of working arount this issue:

* Build a kernel with [the patch that fixes the issue](https://ozlabs.org/~akpm/mmotm/broken-out/hugetlb_cgroup-fix-offline-of-hugetlb-cgroup-with-reservations.patch)
* Disable hugepages in your applications. To do that, remove the hugepage mount and
resource request in your pod deployment file and pass --no-huge to your DPDK app.


## Archive
This is a POC that was built (after significant re-work) based on the work
done for Kubecon 2019. This work can be seen in this repository's history
and in the [archive docs](docs/events/2019-11-KubeCon-NA/README.md)
