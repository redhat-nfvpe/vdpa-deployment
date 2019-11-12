package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang/glog"

	"github.com/openshift/app-netutil/pkg/types"
	netlib "github.com/openshift/app-netutil/lib/v1alpha"
)

var (
	seastarDpdkDynamicFile = flag.String("dpdkfilename", "/var/run/seastar/seastar_dpdk_dynamic.conf", "Name of file to write dpdk data too.")
	seastarIpDynamicFile = flag.String("ipfilename", "/var/run/seastar/seastar_ip_dynamic.conf", "Name of file to write IP data too.")
)

func main() {
	var dpdkInterfaceStr string
	var tmpStr string
	var macStr string
	var ipStr string
	var ipWrtStr string
	var found bool
	var vhostCnt int
	var memifCnt int
	var sriovCnt int

	flag.Parse()

	memifCnt = 1

	//flag.Parse()
	glog.Infof("starting sample application")

	glog.Infof("CALL netlib.GetCPUInfo:")
	cpuResponse, err := netlib.GetCPUInfo()
	if err != nil {
		glog.Errorf("Error calling netlib.GetCPUInfo: %v", err)
		return
	}
	glog.Infof("netlib.GetCPUInfo Response:")
	fmt.Printf("| CPU     |: %+v\n", cpuResponse.CPUSet)

	glog.Infof("CALL netlib.GetInterfaces:")
	ifaceResponse, err := netlib.GetInterfaces()
	if err != nil {
		glog.Errorf("Error calling netlib.GetInterfaces: %v", err)
		return
	}
	glog.Infof("netlib.GetInterfaces Response:")
	for index, iface := range ifaceResponse.Interface {
		tmpStr = ""
		macStr = ""
		ipStr = ""

		fmt.Printf("| %-8d|: IfName=%+v  Name=%+v  Type=%+v\n", index, iface.IfName, iface.Name, iface.Type)
		if iface.Network != nil {
			fmt.Printf("|         |:   %+v\n", iface.Network)
			if iface.Network.Mac != "" {
				macStr = iface.Network.Mac
			}
			for _, ip := range iface.Network.IPs {
				if ip != "" && ipStr == "" {
					ipStr = ip
				}
			}
		}

		switch iface.Type {
		case types.INTERFACE_TYPE_SRIOV:
			if iface.Sriov != nil {
				tmpStr = fmt.Sprintf("-w %s", iface.Sriov.PciAddress)
				sriovCnt++
			}
		case types.INTERFACE_TYPE_VHOST:
			if iface.Vhost != nil {
				if macStr != "" {
					tmpStr = fmt.Sprintf("--vdev virtio_user%d,path=%s,mac=%s,queues=1,queue_size=1024",
						vhostCnt, iface.Vhost.Socketpath, macStr)
				} else {
					tmpStr = fmt.Sprintf("--vdev virtio_user%d,path=%s,queues=1,queue_size=1024",
						vhostCnt, iface.Vhost.Socketpath)
				}
				vhostCnt++
			}
		case types.INTERFACE_TYPE_MEMIF:
			if iface.Memif != nil {
				tmpStr = fmt.Sprintf("--vdev net_memif%d,socket=%s,role=%s%s", memifCnt, iface.Memif.Socketpath, iface.Memif.Role, macStr)
				memifCnt++
			}
		case types.INTERFACE_TYPE_KERNEL, types.INTERFACE_TYPE_VDPA:
		default:
			// For now, do nothing
		}

		if tmpStr != "" {
			if found {
				dpdkInterfaceStr = dpdkInterfaceStr + " "
			}
			dpdkInterfaceStr = dpdkInterfaceStr + tmpStr

			if ipStr != "" && ipWrtStr == "" {
				ipWrtStr = ipStr
			}

			found = true
		}
	}

	if found {
		writeStringToFile(dpdkInterfaceStr, *seastarDpdkDynamicFile)
	}
	if ipWrtStr != "" {
		writeStringToFile(ipWrtStr, *seastarIpDynamicFile)
	}

	return
}

func writeStringToFile(data string, path string) {
	var err error

	fmt.Printf("Writing string: %s\n", data)
	fmt.Printf("Path: %s\n", path)

	dynamicDir := filepath.Dir(path)
	fmt.Printf("Directory: %s\n", dynamicDir)
	if _, err = os.Stat(dynamicDir); err != nil {
		if os.IsNotExist(err) {
			if err = os.MkdirAll(dynamicDir, 0700); err != nil {
				glog.Errorf("Error creating directory: %v", err)
				return
			}
		} else {
			glog.Errorf("Error looking up directory: %v", err)
			return
		}
	}

	file, err := os.Create(path)
	if err != nil {
		glog.Errorf("Error creating file: %v", err)
	} else {
		file.WriteString(data)
		file.Close()
	}
}
