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
	"encoding/json"
	"fmt"
	"runtime"
	_ "flag"

	//v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/containernetworking/cni/pkg/invoke"
	"github.com/containernetworking/cni/pkg/skel"
	cnitypes "github.com/containernetworking/cni/pkg/types"
	"github.com/containernetworking/cni/pkg/types/current"
	cniSpecVersion "github.com/containernetworking/cni/pkg/version"
	"github.com/containernetworking/plugins/pkg/ip"
	"github.com/containernetworking/plugins/pkg/ns"

	_ "github.com/intel/userspace-cni-network-plugin/pkg/k8sclient"
	_ "github.com/intel/userspace-cni-network-plugin/pkg/annotations"
	_ "github.com/intel/userspace-cni-network-plugin/pkg/configdata"

	"github.com/redhat-nfvpe/vdpa-deployment/vdpa-cni/pkg/logging"
	"github.com/redhat-nfvpe/vdpa-deployment/vdpa-cni/pkg/types"
)

var version = "master@git"
var commit = "unknown commit"
var date = "unknown date"

type envArgs struct {
	cnitypes.CommonArgs
	MAC cnitypes.UnmarshallableString `json:"mac,omitempty"`
}

func init() {
	// this ensures that main runs only on main thread (thread group leader).
	// since namespace ops (unshare, setns) are done for a single thread, we
	// must ensure that the goroutine does not jump from OS thread to thread
	runtime.LockOSThread()
}

//
// Local functions
//

func printVersionString() string {
	return fmt.Sprintf("vdpa-cni version:%s, commit:%s, date:%s",
		version, commit, date)
}

// loadNetConf() - Unmarshall the inputdata into the NetConf Structure
func loadNetConf(bytes []byte) (*types.NetConf, error) {
	netconf := &types.NetConf{}
	if err := json.Unmarshal(bytes, netconf); err != nil {
		return nil, fmt.Errorf("loadNetConf(): failed to load netconf: %v", err)
	}

	if netconf.DeviceID == "" {
		return nil, fmt.Errorf("loadNetConf(): VF pci addr is required")
	}

	//
	// Logging
	//
	if netconf.LogFile != "" {
		logging.SetLogFile(netconf.LogFile)
	}
	if netconf.LogLevel != "" {
		logging.SetLogLevel(netconf.LogLevel)
	}

	return netconf, nil
}

func getEnvArgs(envArgsString string) (*envArgs, error) {
	if envArgsString != "" {
		e := envArgs{}
		err := cnitypes.LoadArgs(envArgsString, &e)
		if err != nil {
			return nil, err
		}
		return &e, nil
	}
	return nil, nil
}


func cmdAdd(args *skel.CmdArgs, exec invoke.Exec, kubeClient kubernetes.Interface) error {
	var netConf *types.NetConf

	// Convert the input bytestream into local NetConf structure
	netConf, err := loadNetConf(args.StdinData)

	logging.Infof("cmdAdd: ENTER (AFTER LOAD) - Container %s Iface %s  DeviceID %s",
		args.ContainerID[:12], args.IfName, netConf.DeviceID)
	logging.Verbosef("   Args=%v netConf=%v, exec=%v, kubeClient%v",
		args, netConf, exec, kubeClient)

	if err != nil {
		logging.Errorf("cmdAdd: Parse NetConf - %v", err)
		return err
	}

	envArgs, err := getEnvArgs(args.Args)
	if err != nil {
		logging.Errorf("cmdAdd: Parse args.Args - %v", err)
		return err
	}

	if envArgs != nil {
		MAC := string(envArgs.MAC)
		if MAC != "" {
			netConf.MAC = MAC
		}
	}


	// Initialize returned Result

	// Multus will only copy Interface (i.e. ifName) into NetworkStatus
	// on Pod with Sandbox configured. Get Netns and populate in results.
	netns, err := ns.GetNS(args.Netns)
	if err != nil {
		logging.Errorf("cmdAdd: Failed to open netns %q - %v", netns, err)
		return err
	}
	defer netns.Close()

	result := &current.Result{}
	result.Interfaces = []*current.Interface{{
		Name:    args.IfName,
		Sandbox: netns.Path(),
	}}


	// Retrieve the "SharedDir", directory to create the socketfile in.
	// Save off kubeClient and pod for later use if needed.
	connData, err := getPodAndSharedDir(netConf, args, kubeClient)
	if err != nil {
		logging.Errorf("cmdAdd: Unable to determine \"SharedDir\" - %v", err)
		return err
	}


	saveRemoteConfig(netConf, args, result, connData)

	return cnitypes.PrintResult(result, current.ImplementedSpecVersion)
}


func cmdGet(args *skel.CmdArgs, exec invoke.Exec, kubeClient kubernetes.Interface) error {
/*
	netConf, err := loadNetConf(args.StdinData)

	logging.Infof("cmdGet: (AFTER LOAD) - Container %s Iface %s", args.ContainerID[:12], args.IfName)
	logging.Verbosef("   Args=%v netConf=%v, exec=%v, kubeClient%v",
		args, netConf, exec, kubeClient)

	if err != nil {
		return err
	}

	// FIXME: call all delegates

	return cnitypes.PrintResult(netConf.PrevResult, netConf.CNIVersion)
*/
	return nil
}


func cmdDel(args *skel.CmdArgs, exec invoke.Exec, kubeClient kubernetes.Interface) error {
	var netConf *types.NetConf

	// Convert the input bytestream into local NetConf structure
	netConf, err := loadNetConf(args.StdinData)

	logging.Infof("cmdDel: ENTER (AFTER LOAD) - Container %s Iface %s DeviceID %s", args.ContainerID[:12], args.IfName, netConf.DeviceID)
	logging.Verbosef("   Args=%v netConf=%v, exec=%v, kubeClient%v",
		args, netConf, exec, kubeClient)

	if err != nil {
		logging.Errorf("cmdDel: Parse NetConf - %v", err)
		return err
	}


	// Retrieve the "SharedDir", directory to create the socketfile in.
	// Save off kubeClient and pod for later use if needed.
	/*
	connData, err := getPodAndSharedDir(netConf, args, kubeClient)
	if err != nil {
		logging.Errorf("cmdAdd: Unable to determine \"SharedDir\" - %v", err)
		return err
	}
	*/


	//
	// Cleanup Namespace
	//
	if args.Netns == "" {
		return nil
	}

	err = ns.WithNetNSPath(args.Netns, func(_ ns.NetNS) error {
		var err error
		_, err = ip.DelLinkByNameAddr(args.IfName)
		if err != nil && err == ip.ErrLinkNotFound {
			return nil
		}
		return err
	})

	return err
}

func main() {
	// Init command line flags to clear vendored packages' one, especially in init()
	//flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// add version flag
	//versionOpt := false
	//flag.BoolVar(&versionOpt, "version", false, "Show application version")
	//flag.BoolVar(&versionOpt, "v", false, "Show application version")
	//flag.Parse()
	//if versionOpt == true {
	//	fmt.Printf("%s\n", printVersionString())
	//	return
	//}

	// Extend the cmdAdd(), cmdGet() and cmdDel() functions to take
	// 'exec invoke.Exec' and 'kubeClient k8s.KubeClient' as input
	// parameters. They are passed in as nill from here, but unit test
	// code can then call these functions directly and fake out a
	// Kubernetes Client.
	skel.PluginMain(
		func(args *skel.CmdArgs) error {
			return cmdAdd(args, nil, nil)
		},
		func(args *skel.CmdArgs) error {
			return cmdGet(args, nil, nil)
		},
		func(args *skel.CmdArgs) error {
			return cmdDel(args, nil, nil)
		},
		cniSpecVersion.All,
		"CNI plugin that manages DPDK based interfaces")
}
