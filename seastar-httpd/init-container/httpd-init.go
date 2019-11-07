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
	scyllaDpdkDynamicFile = flag.String("filename", "/var/run/seastar/seastar_dpdk_dynamic.conf", "Name of file to write data too.")
)

func main() {
	var dpdkInterfaceStr string
	var tmpStr string
	var macStr string
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

		fmt.Printf("| %-8d|: IfName=%+v  Name=%+v  Type=%+v\n", index, iface.IfName, iface.Name, iface.Type)
		if iface.Network != nil {
			fmt.Printf("|         |:   %+v\n", iface.Network)
			if iface.Network.Mac != "" {
				macStr = iface.Network.Mac
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
				tmpStr = fmt.Sprintf("--vdev virtio_user%d,path=%s", vhostCnt, iface.Vhost.Socketpath)
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
			found = true
		}
	}

	if found {
		fmt.Printf("Writing string: %s\n", dpdkInterfaceStr)
		fmt.Printf("Path: %s\n", *scyllaDpdkDynamicFile)

		scyllaDpdkDynamicDir := filepath.Dir(*scyllaDpdkDynamicFile)
		fmt.Printf("Directory: %s\n", scyllaDpdkDynamicDir)
		if _, err = os.Stat(scyllaDpdkDynamicDir); err != nil {
			if os.IsNotExist(err) {
				if err := os.MkdirAll(scyllaDpdkDynamicDir, 0700); err != nil {
					glog.Errorf("Error creating directory: %v", err)
					return
				}
			} else {
				glog.Errorf("Error looking up directory: %v", err)
				return
			}
		}

		file, err := os.Create(*scyllaDpdkDynamicFile)
		if err != nil {
			glog.Errorf("Error creating file: %v", err)
		} else {
			file.WriteString(dpdkInterfaceStr)
			file.Close()
		}
	}

	return
}
