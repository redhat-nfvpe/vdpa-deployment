# Sample C Application that calls Net Utility C APIs

## Overview
This directory contains a C program (app_sample.c) that calls the C APIs of the 
GO Net Utility Library. It demonstrates how the APIs can be called and how to
properly free any associatied memory. It also prints out the returned data.

## Quick Start
This section explains an example of building the C sample application that uses
Network Utility Library.

### Compile executable:
To compile the sample C application:
```
$ cd $GOPATH/src/github.com/openshift/app-netutil/
$ make c_sample
```

This builds the GO library as a `.so` file and a C header file under the `bin/`
directory. The sample C application then includes the C header file and the
application binary called `c_sample` is built under `bin/` directory.

### Set LD_LIBRARY_PATH
Before testing, the application needs to know where the shared library is located.
Either copy the `.so` file to a common location (i.e. `/usr/lib/`) or set
`LD_LIBRARY_PATH`:
```
$ echo $LD_LIBRARY_PATH

$ export LD_LIBRARY_PATH=$PWD/bin:$LD_LIBRARY_PATH
```

Note: Printed the original value of `LD_LIBRARY_PATH` so it can be reset later if
desired. Use the following to clear out:
```
$ unset LD_LIBRARY_PATH
```

### Test
Run the application binary:
```
$ ./bin/c_sample
```

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
