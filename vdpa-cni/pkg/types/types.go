// Copyright 2017 Intel Corp.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

import (
	"github.com/containernetworking/cni/pkg/types"
)


type NetConf struct {
	types.NetConf

	DeviceID     string `json:"deviceID"` // PCI address of a VF in valid sysfs format
	MAC          string
	IP           string
	IPMask       string

	// One of the following two must be provided: KubeConfig or SharedDir
	//
	// KubeConfig:
	//  Example: "kubeconfig": "/etc/cni/net.d/multus.d/multus.kubeconfig",
	//  Provides credentials for Userspace CNI to call KubeAPI to:
	//  - Read Volume Mounts:
	//    - "shared-dir": Directory on host socketfiles are created in
	//  - Write annotations:
	//    - "userspace/configuration-data": Configuration data passed
	//      to containe in JSON format.
	//    - "userspace/mapped-dir": Directory in container socketfiles
	//      are created in. Scraped from Volume Mounts above.
	//
	// SharedDir:
	//  Example: "sharedDir": "/usr/local/var/run/openvswitch/",
	//  Since credentials are not provided, Userspace CNI cannot call KubeAPI
	//  to read the Volume Mounts, so this is the same directory used in the
	//  "hostPath". Difference from the "kubeConfig" is that with the "sharedDir"
	//  method, the directory is not unique per POD. That is because the Network
	//  Attachment Definition (where this is defined) is used by multiple PODs.
	//  So this is the base directory and the CNI creates a sub-directory with
	//  the ContainerId as the sub-directory name.
	//
	//  Along the same lines, no annotations are written by Userspace CNI.
	//   1) Configuration data will be written to a file in the same
	//      directory as the socketfiles instead of to an annotation.
	//   2) The "userspace/mapped-dir" annotation must be added to the
	//      pod spec manually (not done by CNI) so container know where to
	//      retrieve data.
	//      Example: userspace/mappedDir: /var/lib/cni/usrspcni/
	KubeConfig    string        `json:"kubeconfig,omitempty"`
	// Not Supported until integrated with new Multus feature
	// of passing config thru CNIAgrs
	SharedDir     string        `json:"sharedDir,omitempty"`

	LogFile       string        `json:"logFile,omitempty"`
	LogLevel      string        `json:"logLevel,omitempty"`
	HostNicFile   string        `json:"hostNic,omitempty"`
}
