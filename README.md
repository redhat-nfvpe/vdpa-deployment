# vdpa-deployment
GO code and example YAML files to deploy vDPA VFs in a container running in kubernetes.

## Overview
vhost Data Path Acceleration (vDPA) utilizes virtio ring compatible
devices to serve virtio driver directly to enable datapath acceleration
(i.e. - vrings are implemented in NIC instead of in software on host).
NICs that support vDPA behave similar to NICs that support SR-IOV in the
fact that the Physical Function (PF) can be divided up into multiple
Virtual Functions (VF). 

This repo, inconjunction with several other repos, enable vDPA VFs to
be used in a container. The following diagram shows the set of components
used and how this repo fits into the end solution:

![](doc/images/DPDKApp_In_Container_Using_vDPA.png)

To leverage this repo, download this repo and follow the steps in the
subsequent sections.

```
   cd $GOPATH/src/
   go get github.com/redhat-nfvpe/vdpa-deployment
   cd github.com/redhat-nfvpe/vdpa-deployment
```

## vdpa-dpdk-image
vDPA leverages existing userspace virtio/vHost protocol to negotiate
the vrings used to pass data traffic. A typical vHost implemenation
uses a userspace vSwitch on the host (like OvS-DPDK or VPP), which
serves as either the server or client of the vHost. Then if running
a VM, QEMU serves as the other side of the vHost (client or server).
If not running a VM but a container, then code like a DPDK application
running in the container serves as the other side of the vHost
(client or server).

For vDPA, where the vrings are being handled by the hardware (NIC),
something needs to handle the vhost negotiation on behalf of the
NIC. In this implementation, this logic is handled by a sample application
provided by the DPDK library (examples/vdpa) which is built and run
in a container as a DaemonSet.

The **vdpa-dpdk-image** directory contains the files to build the
`nfvpe/vdpa-daemonset` docker image. This image runs the vDPA sample
application from DPDK. The `entrypoint.sh` script reads in the set of
PCI Addresses of the vDPA VFs and passes them the the vDPA sample
application. The vDPA sample application then creates the unix socketfiles
and handles the vring negotiation on behalf of the NIC.

Other components in the solution need the PCI Address to socketfile
mapping so the other side of the vhost channel can be setup. The vDPA
sample application has been augmented (with `sed` commands) to export
this mapping to a file (`/var/run/vdpa/socketList.dat`). This file is
read by the gRPC Server (see [server-image](#server-image)), which
exposes the data via a gRPC request/response.

To build the docker image:

```
   cd $GOPATH/src/github.com/redhat-nfvpe/vdpa-deployment
   make vdpa-image
```

This will be deployed later using the **deployment/vdpa-daemonset.yaml**
file.

## server-image
The **server-image** directory contains the files to build the
`nfvpe/vdpa-grpc-server` docker image. This image runs a gRPC Server
that is called from a CNI trying to add a vDPA VF to a container.
In this solution, the SR-IOV CNI has been modified with gRPC Client
code to call this server and retrieve the associated unix socketfile. 

To build the docker image:

```
   cd $GOPATH/src/github.com/redhat-nfvpe/vdpa-deployment
   make server-image
```

This will be deployed later using the **deployment/vdpa-daemonset.yaml**
file.

### client-image
The client code is used only for testing. The testing is not automated
yet. Both the client and server can be built and run locally.

```
   cd $GOPATH/src/github.com/redhat-nfvpe/vdpa-deployment
   make server
   make client
```

Then run the server, which runs until killed:
```
$ sudo ./bin/vdpa-server 
2019/10/04 09:56:30 INFO: loadInterfaces - Using ExampleData.
2019/10/04 09:56:30 INFO: Loaded:
2019/10/04 09:56:30   PCI: 0000:82:00.2  Socketpath: /var/run/vdpa/vhost/vdpa-0
2019/10/04 09:56:30   PCI: 0000:82:00.3  Socketpath: /var/run/vdpa/vhost/vdpa-1
2019/10/04 09:56:30   PCI: 0000:82:00.4  Socketpath: /var/run/vdpa/vhost/vdpa-2
2019/10/04 09:56:30   PCI: 0000:82:00.5  Socketpath: /var/run/vdpa/vhost/vdpa-3
2019/10/04 09:56:30   PCI: 0000:82:00.6  Socketpath: /var/run/vdpa/vhost/vdpa-4
2019/10/04 09:56:30   PCI: 0000:82:00.7  Socketpath: /var/run/vdpa/vhost/vdpa-5
2019/10/04 09:56:30   PCI: 0000:82:01.0  Socketpath: /var/run/vdpa/vhost/vdpa-6
2019/10/04 09:56:30   PCI: 0000:82:01.1  Socketpath: /var/run/vdpa/vhost/vdpa-7
2019/10/04 09:56:30   PCI: 0000:82:01.2  Socketpath: /var/run/vdpa/vhost/vdpa-8
2019/10/04 09:56:30   PCI: 0000:82:01.3  Socketpath: /var/run/vdpa/vhost/vdpa-9
2019/10/04 09:56:30   PCI: 0000:82:01.4  Socketpath: /var/run/vdpa/vhost/vdpa-10
2019/10/04 09:56:30   PCI: 0000:82:01.5  Socketpath: /var/run/vdpa/vhost/vdpa-11
2019/10/04 09:56:30   PCI: 0000:82:01.6  Socketpath: /var/run/vdpa/vhost/vdpa-12
2019/10/04 09:56:30   PCI: 0000:82:01.7  Socketpath: /var/run/vdpa/vhost/vdpa-13
2019/10/04 09:56:30   PCI: 0000:82:02.0  Socketpath: /var/run/vdpa/vhost/vdpa-14
2019/10/04 09:56:30   PCI: 0000:82:02.1  Socketpath: /var/run/vdpa/vhost/vdpa-15
2019/10/04 09:56:30 INFO: Starting vdpaDpdk gRPC Server at: /var/run/vdpa/vdpa.sock
2019/10/04 09:56:30 INFO: vdpaDpdk gRPC Server listening.
```

Then in another window, run the client, which queries the server for
one valid PCI Address, and one unknown PCI Address:
```
$ sudo ./bin/vdpa-client 
2019/10/04 09:56:45 INFO: Starting vDPA-DPDK gRPC Client.
2019/10/04 09:56:45 INFO: Getting socketpath for PCI Address 0000:82:00.5
2019/10/04 09:56:45 INFO: Retrieved socketpath - /var/run/vdpa/vhost/vdpa-3
2019/10/04 09:56:45 INFO: Getting socketpath for PCI Address 0000:87:00.5
2019/10/04 09:56:45 INFO: Retrieved socketpath - 
$
```

The server should then output as the requests are received.
```
2019/10/04 09:56:45 INFO: Received request for PCI Address 0000:82:00.5
2019/10/04 09:56:45 INFO: Found File: /var/run/vdpa/vhost/vdpa-3
2019/10/04 09:56:45 INFO: Received request for PCI Address 0000:87:00.5
2019/10/04 09:56:45 INFO: No File Found!
```

They can also be run in containers
```
   cd $GOPATH/src/github.com/redhat-nfvpe/vdpa-deployment
   make server-image
   make client-image

   kubectl create -f server-image/server.yaml
   kubectl create -f client-image/server.yaml
```

### gRPC proto
Under the **grpc** directory, the file **vdpadpdk-msgs.proto**
defines the gRPC messages used by this server. This file was used
to generate the file **vdpadpdk-msgs.pb.go**. If the messages need
to be updated, then this file needs to be regenerated using:
```
   cd $GOPATH/src/github.com/redhat-nfvpe/vdpa-deployment
   protoc -I grpc/ grpc/vdpadpdk-msgs.proto --go_out=plugins=grpc:grpc/
```

This assumes that `protoc` and `protoc-gen-go` are installed on the
system.

## Deploy the vDPA DaemonSet
The `nfvpe/vdpa-daemonset` and `nfvpe/vdpa-grpc-server` docker images
described above are run in the same pod as a DaemonSet. If the gRPC
Server and Client were running from testing above, make sure they
are stopped.

To deploy (make steps repeated from above):
```
   cd $GOPATH/src/github.com/redhat-nfvpe/vdpa-deployment
   make vdpa-image
   make server-image
   kubectl create -f ./deployment/vdpa-daemonset.yaml
```

**Note:** The above deployment works for Kubernetes 1.16+.
In Kubernetes 1.16, some of the beta features (like DaemonSets)
were removed from beta and promoted to GA. If running an older
version of Kubernetes, use the
`./deployment/pre-1-16-vdpa-daemonset.yaml` file.

The `nfvpe/vdpa-daemonset` image waits for a
`/var/run/vdpa/socketList.dat` to get created by the SR-IOV
Device Plugin. This file is just a list of PCI Address associated
with vDPA VFs. For testing purpose, the file can be generated
manually (where the PCI addresses match the VFs on the server):
```
$ sudo vi /var/run/vdpa/socketList.dat
0000:82:00.2
0000:82:00.3
0000:82:00.4
0000:82:00.5
0000:82:00.6
0000:82:00.7
0000:82:01.0
0000:82:01.1
0000:82:01.2
0000:82:01.3
0000:82:01.4
0000:82:01.5
0000:82:01.6
0000:82:01.7
0000:82:02.0
0000:82:02.1
```
