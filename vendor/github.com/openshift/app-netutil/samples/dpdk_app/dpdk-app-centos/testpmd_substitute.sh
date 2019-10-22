#!/bin/bash

# Add new app-netutil headerfile to the main code so app-netutil
# can be called to gather parameters.
#
# Search for: "#include "testpmd.h"".
# Append:     "#include "dpdk-args.h"".
sed -i -e '/#include "testpmd.h"/a #include "dpdk-args.h"' testpmd.c


# Replace the call to rte_eal_init() to call app-netutil code first
# if no input parametes were passed in. app-netutil code generates
# its own set of DPDK parameters that are used instead. If input
# parameters were passed in, call rte_eal_init() with input parameters
# and run as if app-netutil wasn't there.
#
# Search for the line with "diag = rte_eal_init(argc, argv);"
# Replace that line of code with the contents of file
#   'testpmd_eal_init.txt'.
sed -i '/diag = rte_eal_init(argc, argv);/{
s/diag = rte_eal_init(argc, argv);//g
r testpmd_eal_init.txt
}' testpmd.c


# If no input parametes were passed in, use the parameter list generated
# by app-netutil in the previous patch to call the local parameter
# parsing code, launch_args_parse(). If input parameters were passed in,
# call launch_args_parse() with input parameters and run as if app-netutil
# wasn't there.
#
# Search for the line with "argc -= diag;"
# Create a label 'a' and continue searching and copying until
#   line with "launch_args_parse(argc, argv);" is found.
# Replace that block of code with the contents of file
#   'testpmd_launch_args_parse.txt'.
sed -i '/argc -= diag;/{
:a;N;/launch_args_parse(argc, argv);/!ba;N;s/.*\n//g
r testpmd_launch_args_parse.txt
}' testpmd.c


# Add new app-netutil source file to the Makefile.
#
# Search for line with: "SRCS-y += parameters.c".
# Replace with line:    "SRCS-y += parameters.c dpdk-args.c".
sed -i -e 's/SRCS-y += parameters.c/SRCS-y += parameters.c dpdk-args.c/' Makefile


# Add new app-netutil shared library to the Makefile.
# Contains the C API and GO package which collects the
# interface data.
#
# Search for line with: "SRCS-y += util.c".
# Append line:          "LDLIBS += -lnetutil_api".
sed -i -e '/SRCS-y += util.c/a LDLIBS += -lnetutil_api' Makefile

