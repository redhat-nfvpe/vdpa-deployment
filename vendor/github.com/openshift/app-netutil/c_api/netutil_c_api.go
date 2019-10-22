package main
/*
#include <stdint.h>
#include <stdbool.h>

// Mapped from app-netutil.lib/v1alpha/types.go

struct CPUResponse {
    char*    CPUSet;
};

#define NETUTIL_ERRNO_SUCCESS 0
#define NETUTIL_ERRNO_FAIL 1
#define NETUTIL_ERRNO_SIZE_ERROR 2


#define NETUTIL_NUM_IPS               10
#define NETUTIL_NUM_NETWORKINTERFACE  10

#define NETUTIL_TYPE_UNKNOWN  0
#define NETUTIL_TYPE_KERNEL   1
#define NETUTIL_TYPE_SRIOV    2
#define NETUTIL_TYPE_VHOST    3
#define NETUTIL_TYPE_MEMIF    4
#define NETUTIL_TYPE_VDPA     5

struct NetworkData {
	char*  IPs[NETUTIL_NUM_IPS];
	char*  Mac;
};

struct SriovData {
	char*  PCIAddress;
};


#define NETUTIL_VHOST_MODE_CLIENT  0
#define NETUTIL_VHOST_MODE_SERVER  1
struct VhostData {
	char*  Socketpath;
	int    Mode;
};

#define NETUTIL_MEMIF_ROLE_MASTER  0
#define NETUTIL_MEMIF_ROLE_SLAVE   1
#define NETUTIL_MEMIF_MODE_ETHERNET     0
#define NETUTIL_MEMIF_MODE_IP           1
#define NETUTIL_MEMIF_MODE_INJECT_PUNT  2
struct MemifData {
	char*  Socketpath;
	int    Role;
	int    Mode;
};


struct InterfaceData {
	char*  IfName;
	char*  Name;
	int    Type;
	struct NetworkData Network;
	struct SriovData   Sriov;
	struct VhostData   Vhost;
	struct MemifData   Memif;
};

// *pIface is an array of 'struct InterfaceData' that is allocated
// from the C program.
struct InterfaceResponse {
	int                   numIfaceAllocated;
	int                   numIfacePopulated;
	struct InterfaceData *pIface;
};


*/
import "C"
import "unsafe"

import (
	"flag"

	"github.com/golang/glog"

	netlib "github.com/openshift/app-netutil/lib/v1alpha"
	"github.com/openshift/app-netutil/pkg/types"
)

const (
	cpusetPath = "/sys/fs/cgroup/cpuset/cpuset.cpus"
	netutil_num_ips = 10
	netutil_num_networkstatus = 10
	netutil_num_networkinterface = 10

	// Interface type
	NETUTIL_INTERFACE_TYPE_ALL = types.INTERFACE_TYPE_ALL
	NETUTIL_INTERFACE_TYPE_KERNEL = types.INTERFACE_TYPE_ALL
	NETUTIL_INTERFACE_TYPE_SRIOV = types.INTERFACE_TYPE_SRIOV
	NETUTIL_INTERFACE_TYPE_VHOST = types.INTERFACE_TYPE_VHOST
	NETUTIL_INTERFACE_TYPE_MEMIF = types.INTERFACE_TYPE_MEMIF
	NETUTIL_INTERFACE_TYPE_VDPA = types.INTERFACE_TYPE_VDPA


	// Errno
	NETUTIL_ERRNO_SUCCESS = 0
	NETUTIL_ERRNO_FAIL = 1
	NETUTIL_ERRNO_SIZE_ERROR = 2
)


//export GetCPUInfo
func GetCPUInfo(c_cpuResp *C.struct_CPUResponse) int64 {
	flag.Parse()
	cpuRsp, err := netlib.GetCPUInfo()

	if err == nil {
		c_cpuResp.CPUSet = C.CString(cpuRsp.CPUSet)
		return NETUTIL_ERRNO_SUCCESS
	}
	glog.Errorf("netlib.GetCPUInfo() err: %+v", err)
	return NETUTIL_ERRNO_FAIL
}
//export GetInterfaces
func GetInterfaces(c_ifaceRsp *C.struct_InterfaceResponse) int64 {

	var j C.int

	flag.Parse()
	ifaceRsp, err := netlib.GetInterfaces()

	if err == nil {
		j = 0

		// Map the input pointer to array of structures, c_ifaceResp.pIface, to
		// a slice of the structures, c_ifaceResp_pIface. Then the slice can be
		// indexed.
		c_ifaceResp_pIface := (*[1 << 30]C.struct_InterfaceData)(unsafe.Pointer(c_ifaceRsp.pIface))[:c_ifaceRsp.numIfaceAllocated:c_ifaceRsp.numIfaceAllocated]

		for i, iface := range ifaceRsp.Interface {
			if j < c_ifaceRsp.numIfaceAllocated {
				c_ifaceResp_pIface[j].IfName = C.CString(iface.IfName)
				c_ifaceResp_pIface[j].Name = C.CString(iface.Name)
				c_ifaceRsp.numIfacePopulated++

				if iface.Network != nil {
					c_ifaceResp_pIface[j].Network.Mac =
						C.CString(iface.Network.Mac)

					for k, ip := range iface.Network.IPs {
						if k < netutil_num_ips {
							c_ifaceResp_pIface[j].Network.IPs[k] =
								C.CString(ip)
						} else {
							glog.Errorf("Network.IPs array not sized properly." +
								"At Interface %d, IP index %d.", i, k)
							return NETUTIL_ERRNO_SIZE_ERROR
						}
					}
				}

				switch iface.Type {
				case NETUTIL_INTERFACE_TYPE_KERNEL:
					c_ifaceResp_pIface[j].Type = C.NETUTIL_TYPE_KERNEL
				case NETUTIL_INTERFACE_TYPE_SRIOV:
					c_ifaceResp_pIface[j].Type = C.NETUTIL_TYPE_SRIOV
					if iface.Sriov != nil {
						c_ifaceResp_pIface[j].Sriov.PCIAddress =
							C.CString(iface.Sriov.PciAddress)
					}
				case NETUTIL_INTERFACE_TYPE_VHOST:
					c_ifaceResp_pIface[j].Type = C.NETUTIL_TYPE_VHOST
					if iface.Vhost != nil {
						c_ifaceResp_pIface[j].Vhost.Socketpath =
							C.CString(iface.Vhost.Socketpath)
						if iface.Vhost.Mode == "client" {
							c_ifaceResp_pIface[j].Vhost.Mode = C.NETUTIL_VHOST_MODE_CLIENT
						} else {
							c_ifaceResp_pIface[j].Vhost.Mode = C.NETUTIL_VHOST_MODE_SERVER
						}
					}
				case NETUTIL_INTERFACE_TYPE_MEMIF:
					c_ifaceResp_pIface[j].Type = C.NETUTIL_TYPE_MEMIF
					if iface.Memif != nil {
						c_ifaceResp_pIface[j].Memif.Socketpath =
							C.CString(iface.Memif.Socketpath)
						if iface.Memif.Role == "master" {
							c_ifaceResp_pIface[j].Memif.Role = C.NETUTIL_MEMIF_ROLE_MASTER
						} else {
							c_ifaceResp_pIface[j].Memif.Role = C.NETUTIL_MEMIF_ROLE_SLAVE
						}
						if iface.Memif.Mode == "ethernet" {
							c_ifaceResp_pIface[j].Memif.Mode = C.NETUTIL_MEMIF_MODE_ETHERNET
						} else if iface.Memif.Mode == "ip" {
							c_ifaceResp_pIface[j].Memif.Mode = C.NETUTIL_MEMIF_MODE_IP
						} else {
							c_ifaceResp_pIface[j].Memif.Mode = C.NETUTIL_MEMIF_MODE_INJECT_PUNT
						}
					}
				case NETUTIL_INTERFACE_TYPE_VDPA:
					c_ifaceResp_pIface[j].Type = C.NETUTIL_TYPE_VDPA
				default:
					c_ifaceResp_pIface[j].Type = C.NETUTIL_TYPE_UNKNOWN
				}

				j++
			} else {

				glog.Errorf("InterfaceResponse struct not sized properly." +
					"At Interface %d.", i)

				return NETUTIL_ERRNO_SIZE_ERROR
			}
		}
		return NETUTIL_ERRNO_SUCCESS
	}
	glog.Errorf("netlib.GetInterfaces() err: %+v", err)
	return NETUTIL_ERRNO_FAIL
}

func main() {}
