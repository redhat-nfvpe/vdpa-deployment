#!/bin/bash

# Search for "#include <unistd.h>" and append "#include <stdio.h>".
# Needed for 'fopen' and 'fclose'
sed -i -e '/#include <unistd.h>/a #include <stdio.h>' main.c

# Search for "uint64_t flags;" and append "struct rte_pci_addr addr;".
# This adds the PCI Address to the structure that tracks each device.
sed -i -e '/uint64_t flags;/a struct rte_pci_addr addr;' main.c

# Search for the line with "static int client_mode;" Replace
# that line of code with the contents of 'main_write_table.txt'.
# 'main_write_table.txt' includes "static int client_mode;".
# This adds a function to open a file and write the PCI Address
# and socketfile (in JSON format) in the file.
sed -i '/static int client_mode;/{
s/static int client_mode;//g
r main_write_table.txt
}' main.c

# Search for "vports[devcnt].did = did;" and append
# "vports[devcnt].addr = addr.pci_addr;". In interactive mode,
# when a new device is added to the list, this copies in the PCI Address.
sed -i -e '/vports\[devcnt].did = did;/a vports[devcnt].addr = addr.pci_addr;' main.c

# Search for "vports[i].did = i;" and append the long string.
# In non-interactive mode, when a new deivce is added to the list, this
# retrieves the PCI Address from the DPDK library and adds it to the list.
sed -i -e '/vports\[i].did = i;/a { struct rte_vdpa_device *vdev = rte_vdpa_get_device(i); if (vdev) vports[i].addr = vdev->addr.pci_addr; }' main.c

# Search for "while (scanf("%c", &ch)) {" and insert before "vdpa_write_table();".
# This calls the function inserted above to write the list to a file.
sed -i -e '/while (scanf("%c", &ch)) {/i vdpa_write_table(); while (1) { sleep(9999); }' main.c

