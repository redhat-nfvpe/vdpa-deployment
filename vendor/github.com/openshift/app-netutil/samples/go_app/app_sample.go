package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/golang/glog"

	"github.com/openshift/app-netutil/pkg/types"
	netlib "github.com/openshift/app-netutil/lib/v1alpha"
)

func main() {
	flag.Parse()
	glog.Infof("starting sample application")

	for {
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
			fmt.Printf("| %-8d|: IfName=%+v  Name=%+v  Type=%+v\n", index, iface.IfName, iface.Name, iface.Type)
			if iface.Network != nil {
				fmt.Printf("|         |:   %+v\n", iface.Network)
			}

			switch iface.Type {
			case types.INTERFACE_TYPE_SRIOV:
				if iface.Sriov != nil {
					fmt.Printf("|         |:   PCI=%+v\n", iface.Sriov.PciAddress)
				}
			case types.INTERFACE_TYPE_VHOST:
				if iface.Vhost != nil {
					fmt.Printf("|         |:   Mode=%+v  Socketpath=%+v\n", iface.Vhost.Mode, iface.Vhost.Socketpath)
				}
			case types.INTERFACE_TYPE_MEMIF:
				if iface.Memif != nil {
					fmt.Printf("|         |:   Role=%+v  mode=%+v  Socketpath=%+v\n", iface.Memif.Role, iface.Memif.Mode, iface.Memif.Socketpath)
				}
			case types.INTERFACE_TYPE_KERNEL, types.INTERFACE_TYPE_VDPA:
			default:
				// For now, do nothing
			}
		}
		fmt.Printf("\n")

		time.Sleep(1 * time.Minute)
	}
	return
}
