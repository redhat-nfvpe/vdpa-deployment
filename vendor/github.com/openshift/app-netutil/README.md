# Network Utility Library for Application Running on Kubernetes

## Table of Contents

- [Network Utility](#network-utility)
- [APIs](#apis)
	- [GO APIs](#go-apis)
	- [C APIs](#c-apis)
		- [C Sample APP](#c-sample-app)
		- [DPDK Sample Image](#dpdk-sample-image)
- [Quick Start](#quick-start)
	- [Build GO APP](#build-go-app)
	- [Build C APP](#build-c-app)
	- [Create testpod Image](#create-testpod-image)
	- [Create dpdk-app-centos Image](#create-dpdk-app-centos-image)

## Network Utility
Network Utility (app-netutil) is a library that provides API methods
for applications running in a container to query network information
associated with pod. Network Utility is written in golang and can be
built into an application binary. It also has a C language binding
allowing it to be built into C applications.

To add virtio based interfaces into a DPDK based application in a
container, the DPDK application needs a unix socket file, which is shared
with the host through a VolumeMount, and a set of configuration data about
how the socketfile should be used. Currently, the Userspace CNI uses
annotations or configuration files to share the data between host and
container. SR-IOV needs to get the PCI Addresses of the VFs share with the
DPDK application. Currently it is using Environmental Variables to do this.
Once the above data is in the container, this library has been written to
abstract out where to look and how to process all data passed in.

## APIs
Currently there are two API methods implemented:
* `GetCPUInfo()`
  * This function determines which CPUs are available to the container
  and returns the list to the caller.
* `GetInterfaces()`
  * This function determines the set of interfaces in the container and
  returns the list, along with the interface type and type specific data.

There is a GO and C version of each of these functions.

### GO APIs

There is a GO sample app that provides an example of how to include the
app-netutil as a library in a GO program and how to use the existing APIs:
* [go_app](samples/go_app/README.md)

### C APIs

#### C Sample APP
There is a C sample app that provides an example of how to include the
app-netutil as a library in a C program and how to use the existing APIs:
* [c_app](samples/c_app/README.md)

The gotcha with the C APIs is that data must be allocated on the C side
and passed into the GO library. So the C APIs take data buffers as input.
Then the GO library populates the input structures with the collect data.
It is then up to the C side to free any allocated data.

GO has a special handling for strings where it still allocates the memory
for the string on the C side, but it is hidden in the C.CString() library.
So strings passed from the GO library back to the calling C code must also
free strings, even though there were not explicitly malloc'd. The sample C
code shows examples.

#### DPDK Sample Image
The initial problem `app-netutil` is trying to solve is to collect initial
configuration data for a DPDK application running in a container. The DPDK
Library is written in C, so there is a sample Docker Image that leverages
the C APIs of `app-netutil` to collect the initial configuration data an
then use it to start DPDK. See:
* [dpdk_app](samples/dpdk_app/dpdk-app-centos/README.md)

## Quick Start
This section provides examples of building the sample applications that use
the Network Utility Library. This is just quick start guide, more details
can be found in the links associated with each section.

### Build GO APP

1. Compile executable:
```
$ make go_sample
```
This builds application binary called `go_app` under `bin/` dir.

2. Run:
```
$ ./bin/go_app 
| CPU     |: 0-63
| 0       |: IfName=  Name=cbr0  Type=kernel
|         |:   &{IPs:[10.244.1.32] Mac: Default:true DNS:{Nameservers:[] Domain: Search:[] Options:[]}}
| 1       |: IfName=net1  Name=userspace-ovs-net-1  Type=vhost
|         |:   &{IPs:[10.56.217.166] Mac: Default:false DNS:{Nameservers:[] Domain: Search:[] Options:[]}}
|         |:   Mode=client  Socketpath=/var/lib/cni/usrspcni/3d9fc3b7545d-net1
| 2       |: IfName=net2  Name=userspace-ovs-net-2  Type=vhost
|         |:   &{IPs:[10.77.217.154] Mac: Default:false DNS:{Nameservers:[] Domain: Search:[] Options:[]}}
|         |:   Mode=client  Socketpath=/var/lib/cni/usrspcni/3d9fc3b7545d-net2
| 3       |: IfName=net3  Name=sriov-network  Type=sr-iov
|         |:   &{IPs:[10.56.217.90] Mac:da:18:1d:eb:ef:f2 Default:false DNS:{Nameservers:[] Domain: Search:[] Options:[]}}
| 4       |: IfName=net4  Name=sriov-network  Type=sr-iov
|         |:   &{IPs:[10.56.217.91] Mac:fa:32:38:0e:b5:94 Default:false DNS:{Nameservers:[] Domain: Search:[] Options:[]}}

<CTRL>-C
```

3. Clean up:
```
$ make clean
```

This cleans up built binary and softlinks.


For more details, see:
* [go_app](samples/go_app/README.md)

### Build C APP

1. Compile executable:
```
$ make c_sample
```
This builds application binary called `c_sample` under `bin/`
directory. the `bin\` directory also contains the C header file
`libnetutil_api.h`  and shared library `libnetutil_api.so` needed
to build the C APP.

2. Run:
```
$ export LD_LIBRARY_PATH=$PWD/bin:$LD_LIBRARY_PATH
$ ./bin/c_sample 
Starting sample C application.
Call NetUtil GetCPUInfo():
  cpuRsp.CPUSet = 0-63
Call NetUtil GetInterfaces():
  Interface[0]:
    IfName=""  Name="cbr0"  Type=unknown
    MAC=""  IP="10.244.1.32"
  Interface[1]:
    IfName="net1"  Name="userspace-ovs-net-1"  Type=vHost
    Mode=client  Socketpath="/var/lib/cni/usrspcni/3d9fc3b7545d-net1"
    MAC=""  IP="10.56.217.166"
  Interface[2]:
    IfName="net2"  Name="userspace-ovs-net-2"  Type=vHost
    Mode=client  Socketpath="/var/lib/cni/usrspcni/3d9fc3b7545d-net2"
    MAC=""  IP="10.77.217.154"
  Interface[3]:
    IfName="net3"  Name="sriov-network"  Type=SR-IOV
  
    MAC="da:18:1d:eb:ef:f2"  IP="10.56.217.90"
  Interface[4]:
    IfName="net4"  Name="sriov-network"  Type=SR-IOV
  
    MAC="fa:32:38:0e:b5:94"  IP="10.56.217.91"
```

3. Clean up:
```
$ make clean
```

This cleans up built binary and softlinks.


For more details, see:
* [c_app](samples/c_app/README.md)


### Create testpod Image
The `testpod` image is a CentOS base image built the `app-netutil`
library. It simply creates a container that runs the `go_app` sample
applicatation described above.

1. Build application container image:
```
$ make testpod
```
2. Create application pod:
```
$ kubectl create -f samples/testpod/pod.yaml
```
3. Check for pod logs:
```
$ kubectl logs testpod

I0710 08:07:16.902139       1 app_sample.go:14] starting sample application
I0710 08:07:16.903046       1 resource.go:21] getting cpuset from path: /proc/1/root/sys/fs/cgroup/cpuset/cpuset.cpus
I0710 08:07:16.903574       1 app_sample.go:21] netlib.GetCPUInfo Response: 0-35
I0710 08:07:16.903599       1 resource.go:32] getting environment variables from path: /proc/1/environ
I0710 08:07:16.903669       1 app_sample.go:27] netlib.GetEnv Response:
| KUBERNETES_PORT_443_TCP_PROTO|: tcp
| KUBERNETES_PORT_443_TCP_PORT|: 443
| INSTALL_PKGS             |: golang
| PATH                     |: /usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
| HOSTNAME                 |: testpod
| KUBERNETES_PORT_443_TCP  |: tcp://10.96.0.1:443
| KUBERNETES_SERVICE_HOST  |: 10.96.0.1
| KUBERNETES_SERVICE_PORT  |: 443
| container                |: docker
| HOME                     |: /root
| KUBERNETES_SERVICE_PORT_HTTPS|: 443
| KUBERNETES_PORT          |: tcp://10.96.0.1:443
| KUBERNETES_PORT_443_TCP_ADDR|: 10.96.0.1
I0710 08:07:16.903756       1 network.go:18] getting network status from path: /etc/podinfo/annotations
I0710 08:07:16.904215       1 app_sample.go:36] netlib.GetNetworkStatus Response:
| 0                        |: &{Name: Interface: IPs:[10.96.1.157] Mac:}
...
```
4. Delete application pod:
```
$ kubectl delete -f deployments/pod.yaml
```


For more details, see:
* [go_app](samples/go_app/README.md)

### Create dpdk-app-centos Image
The `dpdk-app-centos` image is a CentOS base image built with DPDK
and includes the `app-netutil` library. The setup to run the image
is more complicated and depends on if you are using vhost interfaces
from something like a Userpace CNI or SR-IOV VFs from SR-IOV CNI.
Below is the quick command to build the image, but it is recommended
that additional README files are consulted for detailed setup
instructions.

1. Build application container image:
```
$ make dpdk_app
```

For more details, see:
* [dpdk_app image](samples/dpdk_app/dpdk-app-centos/README.md)
* [SR-IOV VF Deployment](samples/dpdk_app/sriov/README.md)
