// SPDX-License-Identifier: Apache-2.0
// Copyright(c) 2019 Red Hat, Inc.

//
// Package vdpa manages a vDPA VF. The code retrieves the
// virtio socket file from the vDPA DPDK gRPC Server, then
// writes the relavent virtio/vhost data out so the container
// can negoiate the virtio/vhost connection. The code is leveraging
// the Userspace CNI code to write the data.
//
package vdpa

import (
	"fmt"
	"log"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types/current"

	"github.com/intel/userspace-cni-network-plugin/pkg/annotations"
	"github.com/intel/userspace-cni-network-plugin/pkg/k8sclient"
	usrsptypes "github.com/intel/userspace-cni-network-plugin/pkg/types"

	sriovtypes "github.com/intel/sriov-cni/pkg/types"
)

const (
	defaultKubeConfigFile = "/etc/cni/net.d/multus.d/multus.kubeconfig"
)

// populateConfigData takes the SR-IOV input data and formats ir
// so it can be passed to the Userspace CNI.
func populateConfigData(conf *sriovtypes.NetConf,
	args *skel.CmdArgs,
	ipResult *current.Result) (*usrsptypes.ConfigurationData, error) {
	var err error
	var socketfile string
	var configData usrsptypes.ConfigurationData

	socketfile, err = GetSocketpath(conf.DeviceID)

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

// SaveRemoteConfig writes the input data to annotations (via
// Userspace CNI) for the container to consume.
func SaveRemoteConfig(conf *sriovtypes.NetConf,
	args *skel.CmdArgs,
	ipResult *current.Result) error {

	var kubeClient kubernetes.Interface
	var pod *v1.Pod
	var err error

	kubeConfigFile := defaultKubeConfigFile

	// Retrieve pod so ConfigData can be written to annotations
	pod, kubeClient, err = k8sclient.GetPod(args, kubeClient, kubeConfigFile)
	if err != nil {
		log.Printf("ERROR: SaveRemoteConfig: Failure to retrieve pod - %v", err)
		return err
	}

	// Populate the configData with input data, which will be written to container.
	configData, err := populateConfigData(conf, args, ipResult)
	if err != nil {
		log.Printf("ERROR: SaveRemoteConfig: Failure to retrieve pod - %v", err)
		return err
	}

	// Wrtie configData to the annotations, which will be read by container.
	pod, err = annotations.WritePodAnnotation(kubeClient, pod, configData)
	if err != nil {
		log.Printf("ERROR: SaveRemoteConfig: Failure to write annotations - %v", err)
		return err
	}

	return err
}
