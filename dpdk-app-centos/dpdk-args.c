/* SPDX-License-Identifier: BSD-3-Clause
 * Copyright(c) 2019 Red Hat
 */

#include <stdarg.h>
#include <stdio.h>
#include <stdlib.h>
#include <signal.h>
#include <string.h>
#include <dirent.h>
#include <sys/stat.h>
#include <unistd.h>
#include "libnetutil_api.h"
#include "dpdk-args.h"

bool debugArgs = true;

#define DPDK_ARGS_MAX_ARGS (30)
#define DPDK_ARGS_MAX_ARG_STRLEN (100)
char myArgsArray[DPDK_ARGS_MAX_ARGS][DPDK_ARGS_MAX_ARG_STRLEN];
char* myArgv[DPDK_ARGS_MAX_ARGS];

//#define DPDK_ARGS_MAX_NUM_DIR (30)
//static const char DEFAULT_DIR[] = "/var/lib/cni/";

static char STR_MASTER[] = "master";
static char STR_SLAVE[] = "slave";
static char STR_ETHERNET[] = "ethernet";

/* Large enough to hold: ",mac=aa:bb:cc:dd:ee:ff" */
#define DPDK_ARGS_MAX_MAC_STRLEN (25)
#define DPDK_ARGS_MAX_CONTAINERNAME_STRLEN (80)

static void dumpInterfaces(struct InterfaceResponse *pIfaceRsp) {
	int i, j;
	bool printReturn;

	if ((pIfaceRsp) && (pIfaceRsp->pIface)) {
		for (i = 0; i < pIfaceRsp->numIfacePopulated; i++) {
			printf("  Interface[%d]:\n", i);

			printf("  ");
			printf("  DeviceType=%s",
				(pIfaceRsp->pIface[i].DeviceType == NETUTIL_TYPE_HOST) ? "host" :
				(pIfaceRsp->pIface[i].DeviceType == NETUTIL_TYPE_SRIOV) ? "SR-IOV" :
				(pIfaceRsp->pIface[i].DeviceType == NETUTIL_TYPE_PCI) ? "PCI" :
				(pIfaceRsp->pIface[i].DeviceType == NETUTIL_TYPE_VHOST) ? "vHost" :
				(pIfaceRsp->pIface[i].DeviceType == NETUTIL_TYPE_MEMIF) ? "memif" :
				(pIfaceRsp->pIface[i].DeviceType == NETUTIL_TYPE_VDPA) ? "vDPA" :
				(pIfaceRsp->pIface[i].DeviceType == NETUTIL_TYPE_UNKNOWN) ? "unknown" : "error");

			if (pIfaceRsp->pIface[i].NetworkStatus.Name) {
				printf("  Name=\"%s\"", pIfaceRsp->pIface[i].NetworkStatus.Name);
			}
			if (pIfaceRsp->pIface[i].NetworkStatus.Interface) {
				printf("  Interface=\"%s\"", pIfaceRsp->pIface[i].NetworkStatus.Interface);
			}
			printf("\n");

			printReturn = false;
			if (pIfaceRsp->pIface[i].NetworkStatus.Mac) {
				if (printReturn == false) {
					printReturn = true;
					printf("  ");
				}
				printf("  MAC=\"%s\"", pIfaceRsp->pIface[i].NetworkStatus.Mac);
			}
			for (j = 0; j < NETUTIL_NUM_IPS; j++) {
				if (pIfaceRsp->pIface[i].NetworkStatus.IPs[j]) {
					if (printReturn == false) {
						printReturn = true;
						printf("    DNS Nameservers: ");
					}
					printf("  IP=\"%s\"", pIfaceRsp->pIface[i].NetworkStatus.IPs[j]);
				}
			}
			if (printReturn) {
				printf("\n");
			}

			printReturn = false;
			for (j = 0; j < NETUTIL_NUM_DNS_NAMESERVERS; j++) {
				if (pIfaceRsp->pIface[i].NetworkStatus.DNS.Nameservers[j]) {
					if (printReturn == false) {
						printReturn = true;
						printf("    DNS Nameservers: ");
					}
					printf(" \"%s\"", pIfaceRsp->pIface[i].NetworkStatus.DNS.Nameservers[j]);
				}
			}
			if (printReturn) {
				printf("\n");
			}

			if (pIfaceRsp->pIface[i].NetworkStatus.DNS.Domain) {
				printf("    DNS Domain: \"%s\"\n", pIfaceRsp->pIface[i].NetworkStatus.DNS.Domain);
			}

			printReturn = false;
			for (j = 0; j < NETUTIL_NUM_DNS_SEARCH; j++) {
				if (pIfaceRsp->pIface[i].NetworkStatus.DNS.Search[j]) {
					if (printReturn == false) {
						printReturn = true;
						printf("    DNS Search: ");
					}
					printf(" \"%s\"", pIfaceRsp->pIface[i].NetworkStatus.DNS.Search[j]);
				}
			}
			if (printReturn) {
				printf("\n");
			}

			printReturn = false;
			for (j = 0; j < NETUTIL_NUM_DNS_OPTIONS; j++) {
				if (pIfaceRsp->pIface[i].NetworkStatus.DNS.Options[j]) {
					if (printReturn == false) {
						printReturn = true;
						printf("    DNS Options: ");
					}
					printf(" \"%s\"", pIfaceRsp->pIface[i].NetworkStatus.DNS.Options[j]);
				}
			}
			if (printReturn) {
				printf("\n");
			}

			switch (pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Type) {
				case NETUTIL_TYPE_PCI:
					printf("    Type=PCI");
					if (pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Pci.PciAddress) {
						printf("  PCIAddress=%s", pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Pci.PciAddress);
					}
					if (pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Pci.Vhostnet) {
						printf("  Vhostnet=%s", pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Pci.Vhostnet);
					}
					if (pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Pci.RdmaDevice) {
						printf("  RdmaDevice=%s", pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Pci.RdmaDevice);
					}
					if (pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Pci.PfPciAddress) {
						printf("  PF-PCIAddress=%s", pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Pci.PfPciAddress);
					}
					printf("\n");
					break;
				case NETUTIL_TYPE_VHOST:
					printf("    Type=vHOST");
					printf("  Mode=%s",
						(pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.VhostUser.Mode == NETUTIL_VHOST_MODE_CLIENT) ? "client" :
						(pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.VhostUser.Mode == NETUTIL_VHOST_MODE_SERVER) ? "server" : "error");
					if (pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.VhostUser.Path) {
						printf("  Path=%s", pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.VhostUser.Path);
					}
					printf("\n");
					break;
				case NETUTIL_TYPE_MEMIF:
					printf("    Type=Memif");
					printf("  Role=%s",
						(pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Memif.Role == NETUTIL_MEMIF_ROLE_MASTER) ? "master" :
						(pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Memif.Role == NETUTIL_MEMIF_ROLE_SLAVE) ? "slave" : "error");
					printf("  Mode=%s",
						(pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Memif.Mode == NETUTIL_MEMIF_MODE_ETHERNET) ? "ethernet" :
						(pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Memif.Mode == NETUTIL_MEMIF_MODE_IP) ? "ip" :
						(pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Memif.Mode == NETUTIL_MEMIF_MODE_INJECT_PUNT) ? "inject-punt" : "error");
					if (pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Memif.Path) {
						printf("  Path=%s", pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Memif.Path);
					}
					printf("\n");
					break;
				case NETUTIL_TYPE_VDPA:
					printf("    Type=vDPA");
					if (pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Vdpa.ParentDevice) {
						printf("  ParentDevice=%s", pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Vdpa.ParentDevice);
					}
					if (pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Vdpa.Driver) {
						printf("  Driver=%s", pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Vdpa.Driver);
					}
					if (pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Vdpa.Path) {
						printf("  Path=%s", pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Vdpa.Path);
					}
					if (pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Vdpa.PciAddress) {
						printf("  PCIAddress=%s", pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Vdpa.PciAddress);
					}
					if (pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Vdpa.PfPciAddress) {
						printf("  PF-PCIAddress=%s", pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Vdpa.PfPciAddress);
					}
					printf("\n");
					break;
			}
		}
	}
}

static void freeInterfaces(struct InterfaceResponse *pIfaceRsp) {
	int i, j;

	if ((pIfaceRsp) && (pIfaceRsp->pIface)) {
		for (i = 0; i < pIfaceRsp->numIfacePopulated; i++) {
			if (pIfaceRsp->pIface[i].NetworkStatus.Name) {
				free(pIfaceRsp->pIface[i].NetworkStatus.Name);
			}
			if (pIfaceRsp->pIface[i].NetworkStatus.Interface) {
				free(pIfaceRsp->pIface[i].NetworkStatus.Interface);
			}

			if (pIfaceRsp->pIface[i].NetworkStatus.Mac) {
				free(pIfaceRsp->pIface[i].NetworkStatus.Mac);
			}
			for (j = 0; j < NETUTIL_NUM_IPS; j++) {
				if (pIfaceRsp->pIface[i].NetworkStatus.IPs[j]) {
					free(pIfaceRsp->pIface[i].NetworkStatus.IPs[j]);
				}
			}

			for (j = 0; j < NETUTIL_NUM_DNS_NAMESERVERS; j++) {
				if (pIfaceRsp->pIface[i].NetworkStatus.DNS.Nameservers[j]) {
					free(pIfaceRsp->pIface[i].NetworkStatus.DNS.Nameservers[j]);
				}
			}
			if (pIfaceRsp->pIface[i].NetworkStatus.DNS.Domain) {
				free(pIfaceRsp->pIface[i].NetworkStatus.DNS.Domain);
			}
			for (j = 0; j < NETUTIL_NUM_DNS_SEARCH; j++) {
				if (pIfaceRsp->pIface[i].NetworkStatus.DNS.Search[j]) {
					free(pIfaceRsp->pIface[i].NetworkStatus.DNS.Search[j]);
				}
			}
			for (j = 0; j < NETUTIL_NUM_DNS_OPTIONS; j++) {
				if (pIfaceRsp->pIface[i].NetworkStatus.DNS.Options[j]) {
					free(pIfaceRsp->pIface[i].NetworkStatus.DNS.Options[j]);
				}
			}

			switch (pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Type) {
				case NETUTIL_TYPE_PCI:
					if (pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Pci.PciAddress) {
						free(pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Pci.PciAddress);
					}
					if (pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Pci.Vhostnet) {
						free(pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Pci.Vhostnet);
					}
					if (pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Pci.RdmaDevice) {
						free(pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Pci.RdmaDevice);
					}
					if (pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Pci.PfPciAddress) {
						free(pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Pci.PfPciAddress);
					}
					break;
				case NETUTIL_TYPE_VHOST:
					if (pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.VhostUser.Path) {
						free(pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.VhostUser.Path);
					}
					break;
				case NETUTIL_TYPE_MEMIF:
					if (pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Memif.Path) {
						free(pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Memif.Path);
					}
					break;
				case NETUTIL_TYPE_VDPA:
					if (pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Vdpa.ParentDevice) {
						free(pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Vdpa.ParentDevice);
					}
					if (pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Vdpa.Driver) {
						free(pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Vdpa.Driver);
					}
					if (pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Vdpa.Path) {
						free(pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Vdpa.Path);
					}
					if (pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Vdpa.PciAddress) {
						free(pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Vdpa.PciAddress);
					}
					if (pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Vdpa.PfPciAddress) {
						free(pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Vdpa.PfPciAddress);
					}
					break;
			}
		}

		free(pIfaceRsp->pIface);
	}
}

static int getInterfaces(int argc, int *pPortCnt, int *pPortMask) {
	int i = 0;
	int vhostCnt = 0;
	int memifCnt = 1;
	int sriovCnt = 0;
	int err;
	struct InterfaceResponse ifaceRsp;
	char macStr[DPDK_ARGS_MAX_MAC_STRLEN];
#if 0
	// Refactor code later to find JSON file. Needs work so commented out for now.
	DIR *d;
	struct dirent *dir;
	int currIndex = 0;
	int freeIndex = 0;
	char* dirList[DPDK_ARGS_MAX_NUM_DIR];
	char* fileExt;

	memset(dirList, 0, sizeof(char*)*DPDK_ARGS_MAX_NUM_DIR);

	dirList[freeIndex] = malloc(sizeof(char) * (strlen(DEFAULT_DIR)+1));
	strcpy(dirList[freeIndex++], DEFAULT_DIR);

	while (dirList[currIndex] != NULL) {
		printf("  Directory:%s\n", dirList[currIndex]);
		d = opendir(dirList[currIndex]);
		if (d)
		{
			while ((dir = readdir(d)) != NULL)
			{
				if ((dir->d_name) &&
					(strcmp(dir->d_name, ".") != 0) &&
					(strcmp(dir->d_name, "..") != 0))
				{
					printf("  Name:%s %d\n", dir->d_name, dir->d_type);
					if (dir->d_type == DT_DIR) {
						printf("    Add to Dir List:%s\n", dir->d_name);
						dirList[freeIndex] = malloc(sizeof(char) * (strlen(DEFAULT_DIR)+strlen(dir->d_name)+1));
						sprintf(dirList[freeIndex++], "%s%s/", DEFAULT_DIR, dir->d_name);
					}
					else
					{
						if (strstr(dir->d_name, "net") != NULL)
						{
							fileExt = strrchr(dir->d_name, '.');
							if ((fileExt == NULL) || (strcmp(fileExt, ".json") != 0)) {
								printf("    Adding to vdev list:%s\n", dir->d_name);
								snprintf(&myArgsArray[argc++][0], DPDK_ARGS_MAX_ARG_STRLEN-1,
										 "--vdev=virtio_user%d,path=%s%s", i, dirList[currIndex], dir->d_name);
								i++;
							}
							else {
								printf("    Invalid FileExt\n");
							}
						}
					}
				}
			}
			closedir(d);
		}
		free(dirList[currIndex]);
		currIndex++;
	}
#endif

	ifaceRsp.numIfaceAllocated = NETUTIL_NUM_NETWORKINTERFACE;
	ifaceRsp.numIfacePopulated = 0;
	ifaceRsp.pIface = malloc(ifaceRsp.numIfaceAllocated * sizeof(struct InterfaceData));
	if (ifaceRsp.pIface) {
		memset(ifaceRsp.pIface, 0, (ifaceRsp.numIfaceAllocated * sizeof(struct InterfaceData)));
		err = GetInterfaces(&ifaceRsp);
		if ((err == NETUTIL_ERRNO_SUCCESS) || (err == NETUTIL_ERRNO_SIZE_ERROR)) {

			if (debugArgs) {
				dumpInterfaces(&ifaceRsp);
			}

			for (i = 0; i < ifaceRsp.numIfacePopulated; i++) {
				switch (ifaceRsp.pIface[i].DeviceType) {
					case NETUTIL_TYPE_SRIOV:
						if (ifaceRsp.pIface[i].NetworkStatus.DeviceInfo.Pci.PciAddress) {
							snprintf(&myArgsArray[argc++][0], DPDK_ARGS_MAX_ARG_STRLEN-1,
									 "-w %s", ifaceRsp.pIface[i].NetworkStatus.DeviceInfo.Pci.PciAddress);
							sriovCnt++;

							*pPortMask = *pPortMask | 1 << *pPortCnt;
							*pPortCnt  = *pPortCnt + 1;
						} else {
							printf("ERROR: PCI Address not found. Type=%d\n",
								ifaceRsp.pIface[i].NetworkStatus.DeviceInfo.Type);
						}
						break;
					case NETUTIL_TYPE_VHOST:
						if (ifaceRsp.pIface[i].NetworkStatus.DeviceInfo.VhostUser.Path) {
							if (ifaceRsp.pIface[i].NetworkStatus.DeviceInfo.VhostUser.Mode == NETUTIL_VHOST_MODE_SERVER) {
								snprintf(&myArgsArray[argc++][0], DPDK_ARGS_MAX_ARG_STRLEN-1,
										 "--vdev=virtio_user%d,path=%s,server=1",
										 vhostCnt, ifaceRsp.pIface[i].NetworkStatus.DeviceInfo.VhostUser.Path);

								vhostCnt++;
								*pPortMask = *pPortMask | 1 << *pPortCnt;
								*pPortCnt  = *pPortCnt + 1;
							}
							else if (ifaceRsp.pIface[i].NetworkStatus.DeviceInfo.VhostUser.Mode == NETUTIL_VHOST_MODE_CLIENT) {
								snprintf(&myArgsArray[argc++][0], DPDK_ARGS_MAX_ARG_STRLEN-1,
										 "--vdev=virtio_user%d,path=%s,queues=1",
										 vhostCnt, ifaceRsp.pIface[i].NetworkStatus.DeviceInfo.VhostUser.Path);

								vhostCnt++;
								*pPortMask = *pPortMask | 1 << *pPortCnt;
								*pPortCnt  = *pPortCnt + 1;
							} else {
								printf("ERROR: Unknown vHost Mode=%d\n", ifaceRsp.pIface[i].NetworkStatus.DeviceInfo.VhostUser.Mode);
							}
						} else {
							printf("ERROR: vHost Path not found. Type=%d\n",
								ifaceRsp.pIface[i].NetworkStatus.DeviceInfo.Type);
						}
						break;
					case NETUTIL_TYPE_MEMIF:
						if (ifaceRsp.pIface[i].NetworkStatus.DeviceInfo.Memif.Path) {
							char *pRole = NULL;
							char *pMode = NULL;

							if (ifaceRsp.pIface[i].NetworkStatus.DeviceInfo.Memif.Role == NETUTIL_MEMIF_ROLE_MASTER) {
								pRole = STR_MASTER;
							}
							else if (ifaceRsp.pIface[i].NetworkStatus.DeviceInfo.Memif.Role == NETUTIL_MEMIF_ROLE_SLAVE) {
								pRole = STR_SLAVE;
							}
							else {
								printf("ERROR: Unknown memif Role=%d\n", ifaceRsp.pIface[i].NetworkStatus.DeviceInfo.Memif.Role);
							}

							if (ifaceRsp.pIface[i].NetworkStatus.DeviceInfo.Memif.Mode == NETUTIL_MEMIF_MODE_ETHERNET) {
								pMode = STR_ETHERNET;
							}
							else if (ifaceRsp.pIface[i].NetworkStatus.DeviceInfo.Memif.Mode == NETUTIL_MEMIF_MODE_IP) {
								//pMode = "ip";
								printf("ERROR: memif Mode=%d - Not Supported in DPDK!\n",
									ifaceRsp.pIface[i].NetworkStatus.DeviceInfo.Memif.Mode);
							}
							else if (ifaceRsp.pIface[i].NetworkStatus.DeviceInfo.Memif.Mode == NETUTIL_MEMIF_MODE_INJECT_PUNT) {
								//pMode = "inject-punt"";
								printf("ERROR: memif Mode=%d - Not Supported in DPDK!\n",
									ifaceRsp.pIface[i].NetworkStatus.DeviceInfo.Memif.Mode);
							}
							else {
								printf("ERROR: Unknown memif Mode=%d\n",
									ifaceRsp.pIface[i].NetworkStatus.DeviceInfo.Memif.Mode);
							}

							if ((ifaceRsp.pIface[i].NetworkStatus.Mac) &&
							    (strcmp(ifaceRsp.pIface[i].NetworkStatus.Mac,"") != 0)) {
								snprintf(&macStr[0], DPDK_ARGS_MAX_MAC_STRLEN-1,
										 ",mac=%s", ifaceRsp.pIface[i].NetworkStatus.Mac);
							}
							else {
								macStr[0] = '\0';
							}

							if ((pRole) && (pMode)) {
								snprintf(&myArgsArray[argc++][0], DPDK_ARGS_MAX_ARG_STRLEN-1,
										 "--vdev=net_memif%d,socket=%s,role=%s%s", memifCnt,
										 ifaceRsp.pIface[i].NetworkStatus.DeviceInfo.Memif.Path, pRole, &macStr[0]);

								memifCnt++;
								*pPortMask = *pPortMask | 1 << *pPortCnt;
								*pPortCnt  = *pPortCnt + 1;
							}
						} else {
							printf("ERROR: Memif Path not found. Type=%d\n",
								ifaceRsp.pIface[i].NetworkStatus.DeviceInfo.Type);
						}
						break;
				}
			} /* END of FOR EACH Interface */

			if (sriovCnt == 0) {
				strncpy(&myArgsArray[argc++][0], "--no-pci", DPDK_ARGS_MAX_ARG_STRLEN-1);
			}

			freeInterfaces(&ifaceRsp);
		}
		else {
			printf("Couldn't get network interface, err code: %d\n", err);
		}
	}

	return(argc);
}

static void dumpHugepages(struct HugepagesResponse *pHugepagesRsp) {
	int i;

	if (pHugepagesRsp) {
		if (pHugepagesRsp->MyContainerName) {
			printf("  MyContainerName=%s\n", pHugepagesRsp->MyContainerName);
		}
		if (pHugepagesRsp->pHugepages) {
			for (i = 0; i < pHugepagesRsp->numStructPopulated; i++) {
				printf("  Hugepages[%d]:\n", i);

				printf("  ");
				if (pHugepagesRsp->pHugepages[i].ContainerName) {
					printf("  ContainerName=%s", pHugepagesRsp->pHugepages[i].ContainerName);
				}
				printf("  Request: 1G=%ld 2M=%ld Ukn=%ld  Limit: 1G=%ld 2M=%ld Ukn=%ld\n",
					pHugepagesRsp->pHugepages[i].Request1G,
					pHugepagesRsp->pHugepages[i].Request2M,
					pHugepagesRsp->pHugepages[i].Request,
					pHugepagesRsp->pHugepages[i].Limit1G,
					pHugepagesRsp->pHugepages[i].Limit2M,
					pHugepagesRsp->pHugepages[i].Limit);
			}
		}
	}
}

static void freeHugepages(struct HugepagesResponse *pHugepagesRsp) {
	int i;

	if (pHugepagesRsp) {
		if (pHugepagesRsp->MyContainerName) {
			free(pHugepagesRsp->MyContainerName);
		}
		if (pHugepagesRsp->pHugepages) {
			for (i = 0; i < pHugepagesRsp->numStructPopulated; i++) {
				if (pHugepagesRsp->pHugepages[i].ContainerName) {
					free(pHugepagesRsp->pHugepages[i].ContainerName);
				}
			}
			free(pHugepagesRsp->pHugepages);
		}
	}
}

static int getHugepages(int argc) {
	int i = 0;
	int err;
	int containerIndex = 0;
	int64_t reqMemory = 0;
	int64_t hugepageMemory = 1024;
	struct HugepagesResponse hugepagesRsp;

	memset(&hugepagesRsp, 0, sizeof(struct HugepagesResponse));
	hugepagesRsp.numStructAllocated = NETUTIL_NUM_HUGEPAGES_DATA;
	hugepagesRsp.pHugepages = malloc(hugepagesRsp.numStructAllocated * sizeof(struct HugepagesData));
	if (hugepagesRsp.pHugepages) {
		memset(hugepagesRsp.pHugepages, 0, (hugepagesRsp.numStructAllocated * sizeof(struct HugepagesData)));
		err = GetHugepages(&hugepagesRsp);
		if ((err == NETUTIL_ERRNO_SUCCESS) || (err == NETUTIL_ERRNO_SIZE_ERROR)) {

			if (debugArgs) {
				dumpHugepages(&hugepagesRsp);
			}

			/* Loop through the list of containers to match container name from env. */
			if (hugepagesRsp.MyContainerName) {
				for (i = 0; i < hugepagesRsp.numStructPopulated; i++) {
					if (hugepagesRsp.pHugepages[i].ContainerName) {
						if (strcmp(hugepagesRsp.MyContainerName, hugepagesRsp.pHugepages[i].ContainerName) == 0) {
							containerIndex = i;
							printf("  MATCH: ContainerName=%s, Index=%d\n", hugepagesRsp.pHugepages[i].ContainerName, containerIndex);
							break;
						}
					}
				}
			}

			/* Limit can never be less than Request. So use Limit if non-zero.  */
			/* However, for hugepages, Limit and Request should be the same, so */
			/* either value should be fine.                                     */
			reqMemory =
				(hugepagesRsp.pHugepages[containerIndex].Limit1G != 0) ? hugepagesRsp.pHugepages[containerIndex].Limit1G :
				(hugepagesRsp.pHugepages[containerIndex].Limit2M != 0) ? hugepagesRsp.pHugepages[containerIndex].Limit2M :
				(hugepagesRsp.pHugepages[containerIndex].Limit != 0) ? hugepagesRsp.pHugepages[containerIndex].Limit :
				(hugepagesRsp.pHugepages[containerIndex].Request1G != 0) ? hugepagesRsp.pHugepages[containerIndex].Request1G :
				(hugepagesRsp.pHugepages[containerIndex].Request2M != 0) ? hugepagesRsp.pHugepages[containerIndex].Request2M :
				hugepagesRsp.pHugepages[containerIndex].Request;

			if (reqMemory != 0) {
				/* Assuming 2 NUMA sockets, only use what container has access too. */
				/* TBD: Manage NUMA properly. */ 
				hugepageMemory = reqMemory / 2;
			}

			freeHugepages(&hugepagesRsp);

		} else {
			printf("  Couldn't get Hugepage info, defaulting to %ld, err code: %d\n", hugepageMemory, err);
		}
	}

	/* Build up memory portion of DPDK Args. */
	strncpy(&myArgsArray[argc++][0], "-m", DPDK_ARGS_MAX_ARG_STRLEN-1);
	snprintf(&myArgsArray[argc++][0], DPDK_ARGS_MAX_ARG_STRLEN-1,
		"%ld", hugepageMemory);

	strncpy(&myArgsArray[argc++][0], "-n", DPDK_ARGS_MAX_ARG_STRLEN-1);
	strncpy(&myArgsArray[argc++][0], "4", DPDK_ARGS_MAX_ARG_STRLEN-1);

	return(argc);
}

char** GetArgs(int *pArgc, eDpdkAppType appType)
{
	int argc = 0;
	int i;
	struct CPUResponse cpuRsp;
	int err;
	int portMask = 0;
	int portCnt = 0;
	int lcoreBase = 0;
	int port;
	int length = 0;

	sleep(2);

	memset(&cpuRsp, 0, sizeof(cpuRsp));
	err = GetCPUInfo(&cpuRsp);
	if (err) {
		printf("Couldn't get CPU info, err code: %d\n", err);
	}
	if (cpuRsp.CPUSet) {
		printf("  cpuRsp.CPUSet = %s\n", cpuRsp.CPUSet);

		// Free the string
		free(cpuRsp.CPUSet);
	}


	memset(&myArgsArray[0][0], 0, sizeof(char)*DPDK_ARGS_MAX_ARG_STRLEN*DPDK_ARGS_MAX_ARGS);
	memset(&myArgv[0], 0, sizeof(char)*DPDK_ARGS_MAX_ARGS);

	if (pArgc) {
		/*
		 * Initialize EAL Options
		 */
		strncpy(&myArgsArray[argc++][0], "dpdk-app", DPDK_ARGS_MAX_ARG_STRLEN-1);

		argc = getHugepages(argc);

		//strncpy(&myArgsArray[argc++][0], "--file-prefix=dpdk-app_", DPDK_ARGS_MAX_ARG_STRLEN-1);

		if (appType == DPDK_APP_TESTPMD) {
			strncpy(&myArgsArray[argc++][0], "-l", DPDK_ARGS_MAX_ARG_STRLEN-1);
			strncpy(&myArgsArray[argc++][0], "1-3", DPDK_ARGS_MAX_ARG_STRLEN-1);

			strncpy(&myArgsArray[argc++][0], "--master-lcore", DPDK_ARGS_MAX_ARG_STRLEN-1);
			strncpy(&myArgsArray[argc++][0], "1", DPDK_ARGS_MAX_ARG_STRLEN-1);

			argc = getInterfaces(argc, &portCnt, &portMask);

			/*
			 * Initialize APP Specific Options
			 */
			strncpy(&myArgsArray[argc++][0], "--", DPDK_ARGS_MAX_ARG_STRLEN-1);

			strncpy(&myArgsArray[argc++][0], "--auto-start", DPDK_ARGS_MAX_ARG_STRLEN-1);
			strncpy(&myArgsArray[argc++][0], "--tx-first", DPDK_ARGS_MAX_ARG_STRLEN-1);
			strncpy(&myArgsArray[argc++][0], "--no-lsc-interrupt", DPDK_ARGS_MAX_ARG_STRLEN-1);

			/* testpmd exits if there is not user enteraction, so print stats */
			/* every so often to keep program running. */
			strncpy(&myArgsArray[argc++][0], "--stats-period", DPDK_ARGS_MAX_ARG_STRLEN-1);
			strncpy(&myArgsArray[argc++][0], "60", DPDK_ARGS_MAX_ARG_STRLEN-1);
		}
		else if (appType == DPDK_APP_L3FWD) {
			/* NOTE: The l3fwd app requires a TX Queue per lcore. So seeting lcore to 1 */
			/*       until additional queues are added to underlying interface.         */
			strncpy(&myArgsArray[argc++][0], "-l", DPDK_ARGS_MAX_ARG_STRLEN-1);
			strncpy(&myArgsArray[argc++][0], "1", DPDK_ARGS_MAX_ARG_STRLEN-1);

			strncpy(&myArgsArray[argc++][0], "--master-lcore", DPDK_ARGS_MAX_ARG_STRLEN-1);
			strncpy(&myArgsArray[argc++][0], "1", DPDK_ARGS_MAX_ARG_STRLEN-1);
			lcoreBase = 1;

			argc = getInterfaces(argc, &portCnt, &portMask);

			/*
			 * Initialize APP Specific Options
			 */
			strncpy(&myArgsArray[argc++][0], "--", DPDK_ARGS_MAX_ARG_STRLEN-1);

			/* Set the PortMask, Hexadecimal bitmask of ports used by app. */
			strncpy(&myArgsArray[argc++][0], "-p", DPDK_ARGS_MAX_ARG_STRLEN-1);
			snprintf(&myArgsArray[argc++][0], DPDK_ARGS_MAX_ARG_STRLEN-1,
					"0x%x", portMask);

			/* Set all ports to promiscuous mode so that packets are accepted */
			/* regardless of the packet’s Ethernet MAC destination address.   */
			strncpy(&myArgsArray[argc++][0], "-P", DPDK_ARGS_MAX_ARG_STRLEN-1);

			/* Determines which queues from which ports are mapped to which cores. */
			/* Usage: --config="(port,queue,lcore)[,(port,queue,lcore)]" */
			length = 0;
			for (port = 0; port < portCnt; port++) {
				/* If the first port, add '--config="' to string. */
				if (port == 0) {
					length += snprintf(&myArgsArray[argc][length], DPDK_ARGS_MAX_ARG_STRLEN-length,
									"--config=\"");
				}
				/* If not the first port, add a ',' to string. */
				else {
					length += snprintf(&myArgsArray[argc][length], DPDK_ARGS_MAX_ARG_STRLEN-length, ",");
				}

				/* Add each port data */
				length += snprintf(&myArgsArray[argc][length], DPDK_ARGS_MAX_ARG_STRLEN-length,
					"(%d,%d,%d)", port, 0 /* queue */, lcoreBase /*+port*/);

				/* If the last port, add a trailing " to string. */
				if (port == portCnt-1) {
					length += snprintf(&myArgsArray[argc][length], DPDK_ARGS_MAX_ARG_STRLEN-length, "\"");
				}
			}
			argc++;

			/* Set to use software to analyze packet type. Without this option, */
			/* hardware will check the packet type. Not sure if vHost supports. */
			strncpy(&myArgsArray[argc++][0], "--parse-ptype", DPDK_ARGS_MAX_ARG_STRLEN-1);

		}
		else if (appType == DPDK_APP_L2FWD) {
			strncpy(&myArgsArray[argc++][0], "-l", DPDK_ARGS_MAX_ARG_STRLEN-1);
			strncpy(&myArgsArray[argc++][0], "1-3", DPDK_ARGS_MAX_ARG_STRLEN-1);

			argc = getInterfaces(argc, &portCnt, &portMask);

			/*
			 * Initialize APP Specific Options
			 */
			strncpy(&myArgsArray[argc++][0], "--", DPDK_ARGS_MAX_ARG_STRLEN-1);

			/* Set the PortMask, Hexadecimal bitmask of ports used by app. */
			strncpy(&myArgsArray[argc++][0], "-p", DPDK_ARGS_MAX_ARG_STRLEN-1);
			snprintf(&myArgsArray[argc++][0], DPDK_ARGS_MAX_ARG_STRLEN-1,
					"0x%x", portMask);

			/* Set the PERIOD, statistics will be refreshed each PERIOD seconds. */
			strncpy(&myArgsArray[argc++][0], "-T", DPDK_ARGS_MAX_ARG_STRLEN-1);
			strncpy(&myArgsArray[argc++][0], "120", DPDK_ARGS_MAX_ARG_STRLEN-1);

			/* Set to no-mac-updating. When enabled: */
			/*  - source MAC address is replaced by the TX port MAC address */
		    /*  - The destination MAC address is replaced by 02:00:00:00:00:TX_PORT_ID */
			strncpy(&myArgsArray[argc++][0], "--no-mac-updating", DPDK_ARGS_MAX_ARG_STRLEN-1);
		}

		for (i = 0; i < argc; i++) {
			myArgv[i] = &myArgsArray[i][0];
		}
		*pArgc = argc;
	}

	return(myArgv);
}
