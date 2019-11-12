// SPDX-License-Identifier: Apache-2.0
// Copyright(c) 2019 Red Hat, Inc.

//
// Package types implements the GO types for the project.
//

package hostnic

import (
	"encoding/json"
	"io/ioutil"
	"net"

	"github.com/redhat-nfvpe/vdpa-deployment/vdpa-cni/pkg/logging"
	"github.com/redhat-nfvpe/vdpa-deployment/vdpa-cni/pkg/types"
)

//
// Types
//
type PciMapping struct {
	PciAddress string  `json:"pciAddress"`
	IpAddr string  `json:"ipAddr"`
	IpMask string  `json:"ipMask"`
	Mac string  `json:"mac"`
}

// LoadIpAndMac reads the available vDPA Interfaces (PCI Addresses) from a
// file provided by the network-attachment-definition.
func LoadIpAndMac(netConf *types.NetConf) (bool) {
	var found bool
	var data []byte
	var err error
	
	if netConf.HostNicFile != "" {
		logging.Debugf("LoadIpAndMac - netConf.HostNicFile=%s", netConf.HostNicFile)
		data, err = ioutil.ReadFile(netConf.HostNicFile)
		if err != nil {
			logging.Debugf("LoadIpAndMac - Failed to read input file - %v", err)
			return found
		}

		mappingTable := make([]PciMapping,0)
		if err = json.Unmarshal(data, &mappingTable); err != nil {
			logging.Debugf("LoadIpAndMac - Failed to Unmarshal - %v", err)
			return found
		}

		logging.Debugf("LoadIpAndMac - Loaded: DeviceID=%s", netConf.DeviceID)
		for _, vdpaIface := range mappingTable {
			logging.Debugf("  PCI: %s  IPAddr: %s IPMask: %s MAC: %s",
				vdpaIface.PciAddress, vdpaIface.IpAddr, vdpaIface.IpMask, vdpaIface.Mac)
			if netConf.DeviceID == vdpaIface.PciAddress {
				logging.Debugf("LoadIpAndMac - Found")
				found = true
				netConf.IP = vdpaIface.IpAddr
				netConf.IPMask = vdpaIface.IpMask
				netConf.MAC = vdpaIface.Mac
			}
		}
	} else {
		logging.Debugf("LoadIpAndMac - File not provided.")
	}

	return found
}


// ParseIPv4Mask takes a string and returns an IPMask.
func ParseIPv4Mask(s string) net.IPMask {
	mask := net.ParseIP(s)
	if mask == nil {
		return nil
	}
	return net.IPv4Mask(mask[12], mask[13], mask[14], mask[15])
}

