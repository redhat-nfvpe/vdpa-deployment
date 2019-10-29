// SPDX-License-Identifier: Apache-2.0
// Copyright(c) 2019 Red Hat, Inc.

//
// Package vdpa manages a vDPA VF. The code retrieves the
// virtio socket file from the vDPA DPDK gRPC Server, then
// writes the relavent virtio/vhost data out so the container
// can negoiate the virtio/vhost connection. The code is leveraging
// the Userspace CNI code to write the data.
//
package main

import (
	"fmt"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types/current"

	"github.com/intel/userspace-cni-network-plugin/pkg/annotations"
	"github.com/intel/userspace-cni-network-plugin/pkg/k8sclient"
	usrsptypes "github.com/intel/userspace-cni-network-plugin/pkg/types"

	"github.com/redhat-nfvpe/vdpa-deployment/vdpa-cni/pkg/grpc-client"
	"github.com/redhat-nfvpe/vdpa-deployment/vdpa-cni/pkg/logging"
	"github.com/redhat-nfvpe/vdpa-deployment/vdpa-cni/pkg/types"
)

const (
	defaultKubeConfigFile = "/etc/cni/net.d/multus.d/multus.kubeconfig"
)

type connectionData struct {
	kubeClient kubernetes.Interface
	pod *v1.Pod
	sharedDir string
	kubeConfigFile string
}

func getPodAndSharedDir(netConf *types.NetConf,
						args *skel.CmdArgs,
						kubeClient kubernetes.Interface) (connectionData, error) {

	var connData connectionData
	var found bool
	var err error

	// Retrieve pod so any annotations from the podSpec can be inspected
	connData.pod, connData.kubeClient, err = k8sclient.GetPod(args, kubeClient, netConf.KubeConfig)
	if err != nil {
		logging.Infof("getPodAndSharedDir: Failure to retrieve pod - %v", err)
	} else {

		// Retrieve the sharedDir from the Volumes in podSpec. Directory Socket
		// Files will be written to on host.
		connData.sharedDir, err = annotations.GetPodVolumeMountHostSharedDir(connData.pod)
		if err != nil {
			logging.Infof("getPodAndSharedDir: VolumeMount \"shared-dir\" not provided - %v", err)
		} else {
			found = true
		}
	}

	// Save off the KubeConfig Filename for later use
	connData.kubeConfigFile = netConf.KubeConfig

	err = nil


	if found == false {
		if netConf.SharedDir != "" {
			if netConf.SharedDir[len(netConf.SharedDir)-1:] == "/" {
				connData.sharedDir = fmt.Sprintf("%s%s/", netConf.SharedDir, args.ContainerID[:12])
			} else {
				connData.sharedDir = fmt.Sprintf("%s/%s/", netConf.SharedDir, args.ContainerID[:12])
			}
		} else {
			connData.sharedDir = fmt.Sprintf("%s/%s/", annotations.DefaultBaseCNIDir, args.ContainerID[:12])

			if netConf.KubeConfig == "" {
				logging.Warningf("getPodAndSharedDir: Neither \"KubeConfig\" nor \"SharedDir\" provided, defaulting to %s", connData.sharedDir)
			} else {
				logging.Warningf("getPodAndSharedDir: \"KubeConfig\" invalid and \"SharedDir\" not provided, defaulting to %s", connData.sharedDir)
			}	
		}
	}

	return connData, err
}

// populateUserspaceConfigData takes the SR-IOV input data and formats it
// so it can be passed to the Userspace CNI.
func populateUserspaceConfigData(conf *types.NetConf,
	args *skel.CmdArgs,
	ipResult *current.Result) (*usrsptypes.ConfigurationData, error) {
	var err error
	var socketfile string
	var configData usrsptypes.ConfigurationData

	socketfile, err = vdpagrpcclient.GetSocketpath(conf.DeviceID)

	if socketfile == "" {
		return nil, fmt.Errorf("vDPA: Error looking up socketfile from DeviceID %s", conf.DeviceID)
	}

	//
	// Convert Local Data to usrsptypes.ConfigurationData, which
	// will be written to the container.
	//
	configData.ContainerId = args.ContainerID
	configData.IfName = args.IfName
	configData.Name = conf.Name

	configData.Config.Engine = "vDPA"
	configData.Config.IfType = "vhostuser"
	configData.Config.NetType = "none"

	if ipResult != nil {
		configData.IPResult = *ipResult
	}

	//configData.Config.VhostConf.Mode = "server"
	configData.Config.VhostConf.Mode = "client"
	configData.Config.VhostConf.Socketfile = filepath.Base(socketfile)

	return &configData, err
}

// saveRemoteConfig writes the input data to annotations (via
// Userspace CNI) for the container to consume.
func saveRemoteConfig(conf *types.NetConf,
	args *skel.CmdArgs,
	ipResult *current.Result,
	connData connectionData) error {

	var err error

	// Populate the configData with input data, which will be written to container.
	configData, err := populateUserspaceConfigData(conf, args, ipResult)
	if err != nil {
		logging.Errorf("ERROR: saveRemoteConfig: Failure to retrieve pod - %v", err)
		return err
	}

	// Wrtie configData to the annotations, which will be read by container.
	connData.pod, err = annotations.WritePodAnnotation(connData.kubeClient, connData.pod, configData)
	if err != nil {
		logging.Errorf("ERROR: saveRemoteConfig: Failure to write annotations - %v", err)
		return err
	}

	return err
}
