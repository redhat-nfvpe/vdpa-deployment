#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "libnetutil_api.h"


int main() {
	struct CPUResponse cpuRsp;
	struct InterfaceResponse ifaceRsp;
	int i, j;
	int err;

	printf("Starting sample C application.\n");


	//
	// Example of a C call to GO that returns a string.
	//
	// Note1: Calling C function must free the string.
	//
	printf("Call NetUtil GetCPUInfo():\n");
	memset(&cpuRsp, 0, sizeof(cpuRsp));
	err = GetCPUInfo(&cpuRsp);
	if (err) {
		printf("Couldn't get CPU info, err code: %d\n", err);
		return err;
	}
	if (cpuRsp.CPUSet) {
		printf("  cpuRsp.CPUSet = %s\n", cpuRsp.CPUSet);

		// Free the string
		free(cpuRsp.CPUSet);
	}


	//
	// Example of a C call to GO that returns a structure
	// containing a slice of structures which contains strings.
	//
	// Note1: Calling C function must free the string.
	// Note2: The GO side cannot return any allocated
	//   data, so the data is allocated on the C side and
	//   passed in as a pointer.
	// Note3: Instead of defining the input struct with a fixed
	//   array of entries, the C Program allocates the array
	//   dynamically. For now the number of entries are hardcoded.
	//   Later, could call GO to get the number of entries. 
	//
	printf("Call NetUtil GetInterfaces():\n");
	ifaceRsp.numIfaceAllocated = 10;
	ifaceRsp.numIfacePopulated = 0;
	ifaceRsp.pIface = malloc(ifaceRsp.numIfaceAllocated * sizeof(struct InterfaceData));
	if (ifaceRsp.pIface) {
		memset(ifaceRsp.pIface, 0, (ifaceRsp.numIfaceAllocated * sizeof(struct InterfaceData)));
		err = GetInterfaces(&ifaceRsp);
		if (err) {
			printf("Couldn't get network interface, err code: %d\n", err);
			return err;
		}
		for (i = 0; i < ifaceRsp.numIfacePopulated; i++) {
			printf("  Interface[%d]:\n", i);

			printf("  ");
			if (ifaceRsp.pIface[i].IfName) {
				printf("  IfName=\"%s\"", ifaceRsp.pIface[i].IfName);
				free(ifaceRsp.pIface[i].IfName);
			}
			if (ifaceRsp.pIface[i].Name) {
				printf("  Name=\"%s\"", ifaceRsp.pIface[i].Name);
				free(ifaceRsp.pIface[i].Name);
			}
			printf("  Type=%s",
				(ifaceRsp.pIface[i].Type == NETUTIL_TYPE_KERNEL) ? "kernel" :
				(ifaceRsp.pIface[i].Type == NETUTIL_TYPE_SRIOV) ? "SR-IOV" :
				(ifaceRsp.pIface[i].Type == NETUTIL_TYPE_VHOST) ? "vHost" :
				(ifaceRsp.pIface[i].Type == NETUTIL_TYPE_MEMIF) ? "memif" :
				(ifaceRsp.pIface[i].Type == NETUTIL_TYPE_VDPA) ? "vDPA" :
				(ifaceRsp.pIface[i].Type == NETUTIL_TYPE_UNKNOWN) ? "unknown" : "error");
			printf("\n");

			switch (ifaceRsp.pIface[i].Type) {
				case NETUTIL_TYPE_SRIOV:
					printf("  ");
					if (ifaceRsp.pIface[i].Sriov.PCIAddress) {
						printf("  PCIAddress=%s", ifaceRsp.pIface[i].Sriov.PCIAddress);
						free(ifaceRsp.pIface[i].Sriov.PCIAddress);
					}
					printf("\n");
					break;
				case NETUTIL_TYPE_VHOST:
					printf("  ");
					printf("  Mode=%s",
						(ifaceRsp.pIface[i].Vhost.Mode == NETUTIL_VHOST_MODE_CLIENT) ? "client" :
						(ifaceRsp.pIface[i].Vhost.Mode == NETUTIL_VHOST_MODE_SERVER) ? "server" : "error");
					if (ifaceRsp.pIface[i].Vhost.Socketpath) {
						printf("  Socketpath=\"%s\"", ifaceRsp.pIface[i].Vhost.Socketpath);
						free(ifaceRsp.pIface[i].Vhost.Socketpath);
					}
					printf("\n");
					break;
				case NETUTIL_TYPE_MEMIF:
					printf("  ");
					printf("  Role=%s",
						(ifaceRsp.pIface[i].Memif.Role == NETUTIL_MEMIF_ROLE_MASTER) ? "master" :
						(ifaceRsp.pIface[i].Memif.Role == NETUTIL_MEMIF_ROLE_SLAVE) ? "slave" : "error");
					printf("  Mode=%s",
						(ifaceRsp.pIface[i].Memif.Mode == NETUTIL_MEMIF_MODE_ETHERNET) ? "ethernet" :
						(ifaceRsp.pIface[i].Memif.Mode == NETUTIL_MEMIF_MODE_IP) ? "ip" :
						(ifaceRsp.pIface[i].Memif.Mode == NETUTIL_MEMIF_MODE_INJECT_PUNT) ? "inject-punt" : "error");
					if (ifaceRsp.pIface[i].Memif.Socketpath) {
						printf("  Socketpath=\"%s\"", ifaceRsp.pIface[i].Memif.Socketpath);
						free(ifaceRsp.pIface[i].Memif.Socketpath);
					}
					printf("\n");
					break;
			}

			printf("  ");
			if (ifaceRsp.pIface[i].Network.Mac) {
				printf("  MAC=\"%s\"", ifaceRsp.pIface[i].Network.Mac);
				free(ifaceRsp.pIface[i].Network.Mac);
			}
			for (j = 0; j < NETUTIL_NUM_IPS; j++) {
				if (ifaceRsp.pIface[i].Network.IPs[j]) {
					printf("  IP=\"%s\"", ifaceRsp.pIface[i].Network.IPs[j]);
					free(ifaceRsp.pIface[i].Network.IPs[j]);
				}
			}
			printf("\n");
		}
	}

	return 0;
}
