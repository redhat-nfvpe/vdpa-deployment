# KubeCon - North America: San Diego, CA  November 26-28, 2019
Red Hat presented two demos at the Red Hat Booth at KubeCon NA
in San Diego, on November 26-28, 2019.
* Breaking Cloud Native Network Performance Barriers
* Making High Performance Networking Applications Work On Hybrid Clouds

## Breaking Cloud Native Network Performance Barriers
With the rise of Kubernetes, many legacy VNFs are being migrated
over to CNFs. As 5G comes closer to a reality, Service Providers
and Telcos are moving their workloads to containers and adding
new requirements to container networking. These requirements
include higher bandwidth and lower latency packet processing as
well as passing Layer 2 and Layer 3 traffic into the container
(as opposed to typical container expecting only Layer 7 traffic).

Using vDPA, along with technolgies like Multus to inject additional
interfaces into a container, these new requirements can be met. This
presentation discussed the components used to build such a deployment,
as well as presenting a live demonstration of a POC with all the
components actually deployed.
[Slides](2019-11-KubeConNA-BreakingCloudNativeNetworkPerformanceBarriers.pdf)

Multus was used inject additional interfaces into a container.
[Multus](https://github.com/intel/multus-cni) is a meta-plugin. It is the
only CNI called from Kublet. Multus is provided with a set of delgates
(CNIs to call) and their associated configuration data, and Multus cycles
through each delgate, calling each one and saving the results. Then the
result of the 'default' interface is returned to Kubelet. Kubernetes only
know about the default interfaces, eth0, and nothing about the additional
interfaces (net0-netn). Today, Multus is being used to inject SR-IOV
interfaces into containers.

Virtual Data Path Acceleration (vDPA) utilizes virtio ring compatible
devices to serve virtio driver directly to enable datapath acceleration
(i.e. - virtio vrings are implemented in NIC instead of in software on
host). Similar to the virtio full HW offloading the data plane uses a
standard virtio device with standard virtio ring layout. The data plane
goes directly from the NIC to the workload using these virtio rings.
However each NIC vendor can now continue using its own driver (with a
small vDPA add-on) and a generic vDPA driver is added to the kernel to
translate the vendor NIC driver/control-plane to the virtio control
plane.

![](images/BreakingCloudNativeNetworkPerformanceBarriers_BM.png)

Components of POC:
* **vDPA Device Plugin**: This is a slightly modified version of the
  SR-IOV Device Plugin. Given a selection criteria (via a configMap),
  the vDPA Device Plugin, running as a Container DaemonSet, detects a
  set of usable PCI Addresses and adds them to a resource pool. When
  the workload container is spun up, it requests resources from this
  pool. The PCI Address is then reserved and passed through Kubelet to
  CNIs (**Multus** then to **vDPA CNI**).
* **DPDK vDPA**: This container contains the vendor specific code. Given
  a set of PCI Addresses (from the **vDPA Device Plugin** to
  init-container via configuration file), it creates a vHost socket file
  for each PCI Address. This socketfile is passed into the workload
  container (**DPDK l3fwd**) and the virtio control plane is negotiated
  through this socket file. Then the vendor specific code manages the
  physical NIC.
* **gRPC Server**: This container is a simple gRPC server, which takes
  in a PCI Address in the request and responds with the associated socket
  file.
* **Multus**: As described above, **Multus** is used to inject additional
  interfaces into the container. The 'default' interface is managed by
  Flannel, then two additional vDPA interfaces are injected.
* **vDPA CNI**: The **vDPA CNI** is used to passed the vHost socket files
  into the container so the virtio control plane can negogiate the data
  plane vrings.
* **DPDK l3fwd**: This is the workload container. The DPDK sample
  application, l3fwd, is being run in the container. DPDK has the vhost
  protocol built in, so it can negogiate and manage the virtio devices,
  which includes a vDPA device.

For more details on all the components used in the PoC, see [Details](DETAILS.md)

## Making High Performance Networking Applications Work On Hybrid Clouds
This presentation builds on the previous presentation. Where as
the first concentrates on using a high-speed vDPA interface in a
container workload, this presentation focuses on the fact that
because vDPA is using an open and standard protocol, the vendor
specific logic is outside the container workload. Therefore, the
same container image can be run across multiple platforms with
different NICs, some supporting vDPA and others not. During this
presentation, the same contianer image was run on Bare Metal, on
Alibaba Bare Metal Cloud, and on AWS Cloud.
[Slides](2019-11-KubeConNA-MakingHighPerformanceNetworkingApplicationsWorkOnHybridClouds.pdf)

In all three deployments, the workload is the
[Seastar httpd](https://github.com/scylladb/seastar) from Scylla.
This package comes with an HTTP loader application, seawreck, which
will be used to drive httpd.

### Bare Metal
![](images/MakingHighPerformanceNetworkingApplicationsWorkOnHybridClouds_BM.png)

Components of POC:
* **vDPA Device Plugin**: This is a slightly modified version of the
  SR-IOV Device Plugin. Given a selection criteria (via a configMap),
  the vDPA Device Plugin, running as a Container DaemonSet, detects a
  set of usable PCI Addresses and adds them to a resource pool. When
  the workload container is spun up, it requests resources from this
  pool. The PCI Address is then reserved and passed through Kubelet to
  CNIs (**Multus** then to **vDPA CNI**).
* **DPDK vDPA**: This container contains the vendor specific code. Given
  a set of PCI Addresses (from the **vDPA Device Plugin** to
  init-container via configuration file), it creates a vHost socket file
  for each PCI Address. This socketfile is passed into the workload
  container (**DPDK l3fwd**) and the virtio control plane is negotiated
  through this socket file. Then the vendor specific code manages the
  physical NIC.
* **gRPC Server**: This container is a simple gRPC server, which takes
  in a PCI Address in the request and responds with the associated socket
  file.
* **Multus**: As described above, **Multus** is used to inject additional
  interfaces into the container. The 'default' interface is managed by
  Flannel, then an additional vDPA interface is injected.
* **vDPA CNI**: The **vDPA CNI** is used to passed the vHost socket files
  into the container so the virtio control plane can negogiate the data
  plane vrings.
* **Seastar httpd**: This is the workload container. This is a DPDK
  based application being run in the container. DPDK has the vhost
  protocol built in, so it can negogiate and manage the virtio devices,
  which includes a vDPA device.

### Alibaba Bare Metal Cloud
![](images/MakingHighPerformanceNetworkingApplicationsWorkOnHybridClouds_Alibaba.png)

Alibaba Bare Metal Cloud does not support any NICs that support
vDPA. However, the Alibaba Bare Metal Cloud does support virtio
full HW Offloading. Instead of mapping VFs from a vDPA PF NIC,
two virtio full HW Offloaded NICs were used. This required a
different set of selectors in configMap passed to the **vDPA
Device Plugin** (configuration change, not code change).

Components of POC:
* **vDPA Device Plugin**: Same as above.
* **DPDK vDPA**: This container contains the vendor specific code.
  Some additional patches were needed to support the virtio full
  HW Offloading.
* **gRPC Server**: Same as above.
* **Multus**: Same as above.
* **vDPA CNI**: Same as above.
* **Seastar httpd**: Same as above.

### AWS Cloud
![](images/MakingHighPerformanceNetworkingApplicationsWorkOnHybridClouds_AWS.png)

AWS Cloud does not support any NICs that support vDPA or
virtio full HW Offloading. So in the AWS deployment, a
translation layer was added. **DPDK testpmd** adds a
translation from the AWS support NIC to a vHost interface.

Components of POC:
* **vDPA Device Plugin**: Same as above.
* **DPDK vDPA**: Removed and replaced by **DPDK testpmd**.
* **DPDK testpmd**: This layer adds a translation from the
  AWS support NIC to a vHost interface. This introduces a
  packet copy and lower performance, but the same workload
  was able to be run as is.
* **gRPC Server**: Same as above.
* **Multus**: Same as above.
* **vDPA CNI**: Same as above.
* **Seastar httpd**: Same as above.

## Packages
The following repo versions were used in the demos presented
at the Red Hat Booth at KubeCon NA in San Diego, on November
26-28, 2019:
* https://github.com/redhat-nfvpe/vdpa-deployment
  * Commit: a8732f7c351d79d7977dcebcfb5e75dce74611f7 from 11/18/2019
    * Commit demo actually ran on.
    * Tag '2019-11-KubeCon-NA' added after the fact to include
      slides and lock git requests to particular commit used in
      the demo.
* https://github.com/intel/sriov-network-device-plugin
  * Commit: bf28fdc3e2d9dd2edcc4f0bceb58448ac9317696 from 10/17/2019
  * Image: sriov-dp
  * https://github.com/intel/userspace-cni-network-plugin
  * Commit: 323d722f1046b4224002c10b09dd529b1432c1d6 from 11/04/2019
  * Image: vDPA CNI
* https://github.com/openshift/app-netutil
  * Commit: 156a6eab4812fbfaa3173ace626a453b440da7db from 10/22/2019
  * Image: dpdk-app-centos
  * Image: seastart-httpd/init-container
* https://gitlab.com/mcoquelin/dpdk-next-virtio
  * Branch: kubecon_vdpa_drivers
    * Commit: 32bf71c3b6774faa0029a62bf6ee9334a7980c84 from 11/13/2019
    * Image: dpdk-app-centos
    * Image: vdpa-dpdk-image
  * Branch: seastar_kubecon_poc
    * Commit: 0c8ef461ba4525bb2b6a5636cb79a5a7d10df75c from 11/13/2019
    * Image: seastart-httpd/httpd
* https://gitlab.com/mcoquelin/seastar
  * Branch: kubecon_poc_timer_workaround
    * Commit: 7a454cc105dc24a97fa06b10124cf3850d77d4e9 from 11/15/2019
    * Image: seastart-httpd/httpd
