// SPDX-License-Identifier: Apache-2.0
// Copyright(c) 2019 Red Hat, Inc.

//
// This module implements the server side for the gRPC calls.
// This code is intended run in a sidecar to the vDPA-DPDK
// application. The vDPA-DPDK application serves as one side,
// the hardware vring side, of the vHost channel that is used
// by the vDPA interfaces to negotiate with vring settings.
//
// The gRPC Serber provides an API for CNIs to retrieve the
// Unix Socket file of the vHost from the vDPA-DPDK application.
//

package main

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang/glog"

	"google.golang.org/grpc"

	vdpagrpc "github.com/redhat-nfvpe/vdpa-deployment/grpc"
	vdpatypes "github.com/redhat-nfvpe/vdpa-deployment/pkg/types"
)

var (
	jsonDBFile = flag.String("json_db_file", "", "A json file containing a list of vDPA VF Interfaces")
)

type vdpaDpdkServer struct {
	grpcServer *grpc.Server
	mappingTable []*vdpatypes.VFPCIMapping
}

// GetSocketpath takes a PCI Address of a vDPA VF and returns the
// associated Socketpath
func (s *vdpaDpdkServer) GetSocketpath(ctx context.Context, req *vdpagrpc.GetSocketpathRequest) (*vdpagrpc.GetSocketpathResponse, error) {
	var rsp vdpagrpc.GetSocketpathResponse

	glog.Infof("Received request for PCI Address %s", req.PciAddress)
	for _, vdpaIface := range s.mappingTable {
		if vdpaIface.PciAddress == req.PciAddress {
			rsp.Socketpath = vdpaIface.Socketpath
			glog.Infof("Found File: %s", rsp.Socketpath)
			break
		}
	}
	if rsp.Socketpath == "" {
		glog.Infof("No File Found!")
	}
	
	return &rsp, nil
}

// loadInterfaces reads the available vDPA Interfaces (PCI Addresses) from a
// file provided by the init container.
func (s *vdpaDpdkServer) loadInterfaces(filePath string) {
	var data []byte
	if filePath != "" {
		var err error
		data, err = ioutil.ReadFile(filePath)
		if err != nil {
			glog.Errorf("Failed to read input file: %v", err)
		}
	} else {
		data = exampleData
	}
	if err := json.Unmarshal(data, &s.mappingTable); err != nil {
		glog.Errorf("Failed to read sample data: %v", err)
	}
}

func (s *vdpaDpdkServer) start() error {
	glog.Infof("starting vdpaDpdk server at: %s\n", vdpatypes.GRPCEndpoint)
	lis, err := net.Listen("unix", vdpatypes.GRPCEndpoint)
	if err != nil {
		glog.Errorf("Error creating vdpaDpdk gRPC service: %v", err)
		return err
	}

	vdpagrpc.RegisterVdpaDpdkServer(s.grpcServer, s)
	go s.grpcServer.Serve(lis)

	// Wait for server to start
	conn, err := grpc.Dial(vdpatypes.GRPCEndpoint, grpc.WithInsecure(), grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
		grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("unix", addr, timeout)
		}),
	)
	if err != nil {
		glog.Errorf("Error starting vdpaDpdk server: %v", err)
		return err
	}
	glog.Infof("vdpaDpdk server start serving")
	conn.Close()
	return nil
}

func (s *vdpaDpdkServer) stop() error {
	glog.Infof("stopping vdpaDpdk server")
	if s.grpcServer != nil {
		s.grpcServer.Stop()
		s.grpcServer = nil
	}
	err := os.Remove(vdpatypes.GRPCEndpoint)
	if err != nil && !os.IsNotExist(err) {
		glog.Errorf("Error cleaning up socket file")
	}
	return nil
}

// newVdpaDpdkServer creates a new server instance with initialized data.
func newVdpaDpdkServer() *vdpaDpdkServer {
	s := &vdpaDpdkServer{
		grpcServer: grpc.NewServer(),
	}
	s.loadInterfaces(*jsonDBFile)
	return s
}

func main() {
	flag.Parse()

	// Cleanup socketfile before starting.
	err := os.Remove(vdpatypes.GRPCEndpoint)
	if err != nil && !os.IsNotExist(err) {
		glog.Errorf("Error cleaning up socket file")
	}

	vdpaServer := newVdpaDpdkServer()
	if vdpaServer == nil {
		glog.Errorf("Error initializing netutil manager")
		return
	}
	err = vdpaServer.start()
	if err != nil {
		return
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case sig := <-sigCh:
		glog.Infof("signal received, shutting down", sig)
		vdpaServer.stop()
		return
}}


// exampleData is a copy of vdpa_db.json. It's to avoid
// specifying file path with `go run`.
var exampleData = []byte(`[{
    "pciAddress": "0000:82:00.2",
    "socketpath": "/var/run/vdpa/vdpa-0"
}, {
    "pciAddress": "0000:82:00.3",
    "socketpath": "/var/run/vdpa/vdpa-1"
}, {
    "pciAddress": "0000:82:00.4",
    "socketpath": "/var/run/vdpa/vdpa-2"
}, {
    "pciAddress": "0000:82:00.5",
    "socketpath": "/var/run/vdpa/vdpa-3"
}, {
    "pciAddress": "0000:82:00.6",
    "socketpath": "/var/run/vdpa/vdpa-4"
}, {
    "pciAddress": "0000:82:00.7",
    "socketpath": "/var/run/vdpa/vdpa-5"
}, {
    "pciAddress": "0000:82:01.0",
    "socketpath": "/var/run/vdpa/vdpa-6"
}, {
    "pciAddress": "0000:82:01.1",
    "socketpath": "/var/run/vdpa/vdpa-7"
}, {
    "pciAddress": "0000:82:01.2",
    "socketpath": "/var/run/vdpa/vdpa-8"
}, {
    "pciAddress": "0000:82:01.3",
    "socketpath": "/var/run/vdpa/vdpa-9"
}, {
    "pciAddress": "0000:82:01.4",
    "socketpath": "/var/run/vdpa/vdpa-10"
}, {
    "pciAddress": "0000:82:01.5",
    "socketpath": "/var/run/vdpa/vdpa-11"
}, {
    "pciAddress": "0000:82:01.6",
    "socketpath": "/var/run/vdpa/vdpa-12"
}, {
    "pciAddress": "0000:82:01.7",
    "socketpath": "/var/run/vdpa/vdpa-13"
}, {
    "pciAddress": "0000:82:02.0",
    "socketpath": "/var/run/vdpa/vdpa-14"
}, {
    "pciAddress": "0000:82:02.1",
    "socketpath": "/var/run/vdpa/vdpa-15"
}]`)
