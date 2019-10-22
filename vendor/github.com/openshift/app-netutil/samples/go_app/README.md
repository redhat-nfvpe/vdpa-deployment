# Sample GO Application that calls Net Utility GO APIs

## Overview
This directory contains a GO program (app_sample.go) that calls the 
GO Net Utility Library. It demonstrates how the APIs can be called.
It also prints out the returned data.

## Quick Start
This section explains an example of building the GO sample application
that uses Network Utility Library.

### Compile executable
To compile the sample GO application:
```
$ cd $GOPATH/src/github.com/openshift/app-netutil/
$ make
```

This builds the GO sample application binary called `go_app` under the
`bin/` directory.

### Test
Run the application binary:
```
$ ./bin/go_app
```

This application run a forever loop, calling the Network Utility Library
and printing the results. Then slepping for 1 minute and repeating. Use
<CTRL>-C to exit.

#### Debug Logs
To see additional debug messages, pass logging information as input to the
sample GO application:
```
$ ./bin/go_app -stderrthreshold=INFO
```

Valid log levels are:
* ERROR
* WARNING
* INFO

#### Run locally
If the application is not actually running in a container where annotations have been
exposed, run the following to copy a sample annotation file onto the system. There are
a couple of examples, so choose one that suits your testing. Make sure to name the
file 'annotations' in the '/etc/podinfo/' directory.
```
$ sudo mkdir -p /etc/podinfo/
$ sudo cp samples/annotations/annotations_all /etc/podinfo/annotations
```

SR-IOV exposes the PCI Addresses of the VF to the container using an
environmental variable. If the application is not actually running in a
container where the SR-IOV environmental variables have been created, run
the following to set a sample environmental variable in the session:
```
$ export PCIDEVICE_INTEL_COM_SRIOV=0000:03:02.1,0000:03:04.3
$ echo $PCIDEVICE_INTEL_COM_SRIOV
0000:03:02.1,0000:03:04.3
```

### Clean up
To cleanup all generated files, run:
```
$ make clean
```

This cleans up built binary and other generated files.
