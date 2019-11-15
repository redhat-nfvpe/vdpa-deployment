// SPDX-License-Identifier: Apache-2.0
// Copyright(c) 2019 Red Hat, Inc.

//
// Package main implements the server side for the gRPC calls.
// This code is intended run in a sidecar to the vDPA-DPDK
// application. The vDPA-DPDK application serves as one side,
// the hardware vring side, of the vHost channel that is used
// by the vDPA interfaces to negotiate with vring settings.
//
// The gRPC Server provides an API for CNIs to retrieve the
// Unix Socket file of the vHost from the vDPA-DPDK application.
//

package main

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"

	vdpagrpc "github.com/redhat-nfvpe/vdpa-deployment/grpc"
	vdpatypes "github.com/redhat-nfvpe/vdpa-deployment/pkg/types"
)

var (
	socketlist = flag.String("socketlist", "", "A json file containing a list of vDPA VF Interfaces")
)

type vdpaDpdkServer struct {
	grpcServer *grpc.Server
	mappingTable []*vdpatypes.VFPCIMapping
}

// GetSocketpath takes a PCI Address of a vDPA VF and returns the
// associated Socketpath
func (s *vdpaDpdkServer) GetSocketpath(ctx context.Context, req *vdpagrpc.GetSocketpathRequest) (*vdpagrpc.GetSocketpathResponse, error) {
	var rsp vdpagrpc.GetSocketpathResponse

	log.Printf("INFO: Received request for PCI Address %s", req.PciAddress)
	for _, vdpaIface := range s.mappingTable {
		if vdpaIface.PciAddress == req.PciAddress {
			rsp.Socketpath = vdpaIface.Socketpath
			log.Printf("INFO: Found File: %s", rsp.Socketpath)
			break
		}
	}
	if rsp.Socketpath == "" {
		log.Printf("INFO: No File Found!")
	}
	
	return &rsp, nil
}

// loadInterfaces reads the available vDPA Interfaces (PCI Addresses) from a
// file provided by the init container.
func (s *vdpaDpdkServer) loadInterfaces(filePath string) {
	var data []byte
	if filePath != "" {
		var err error
		log.Printf("INFO: loadInterfaces - filePath=%s\n", filePath)
		data, err = ioutil.ReadFile(filePath)
		if err != nil {
			log.Printf("ERROR: Failed to read input file - %v", err)
		}
		/*
		err = os.Remove(filePath)
		if err != nil {
			log.Printf("ERROR: Failed to remove file - %v", err)
		}
		*/
	} else {
		log.Printf("INFO: loadInterfaces - Using ExampleData.\n")
		data = exampleData
	}
	if err := json.Unmarshal(data, &s.mappingTable); err != nil {
		log.Printf("ERROR: Failed to read ExampleData - %v", err)
	}

	log.Printf("INFO: Loaded:")
	for _, vdpaIface := range s.mappingTable {
		log.Printf("  PCI: %s  Socketpath: %s", vdpaIface.PciAddress, vdpaIface.Socketpath)
	}
}

func (s *vdpaDpdkServer) start() error {
	log.Printf("INFO: Starting vdpaDpdk gRPC Server at: %s\n", vdpatypes.GRPCEndpoint)
	lis, err := net.Listen("unix", vdpatypes.GRPCEndpoint)
	if err != nil {
		log.Printf("ERROR: Error opening socket - %v", err)
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
		log.Printf("ERROR: Error dialing socket - %v", err)
		return err
	}
	log.Printf("INFO: vdpaDpdk gRPC Server listening.")
	conn.Close()
	return nil
}

func (s *vdpaDpdkServer) stop() error {
	log.Printf("INFO: Stopping vdpaDpdk gRPC Server.")
	if s.grpcServer != nil {
		s.grpcServer.Stop()
		s.grpcServer = nil
	}
	err := os.Remove(vdpatypes.GRPCEndpoint)
	if err != nil && !os.IsNotExist(err) {
		log.Printf("ERROR: Error cleaning up socket file.")
	}
	return nil
}

// newVdpaDpdkServer creates a new server instance with initialized data.
func newVdpaDpdkServer() *vdpaDpdkServer {
	s := &vdpaDpdkServer{
		grpcServer: grpc.NewServer(),
	}
	s.loadInterfaces(*socketlist)
	return s
}

func main() {
	flag.Parse()

	// Cleanup socketfile before starting.
	err := os.Remove(vdpatypes.GRPCEndpoint)
	if err != nil && !os.IsNotExist(err) {
		log.Printf("INFO: Error cleaning up previous socket file.")
	}

	vdpaServer := newVdpaDpdkServer()
	if vdpaServer == nil {
		log.Printf("ERROR: Error creating new vdpaDpdk gRPC Server.")
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
		log.Printf("INFO: signal received, shutting down.", sig)
		vdpaServer.stop()
		return
}}


// exampleData is a copy of vdpa_db.json. It's to avoid
// specifying file path with `go run`.
var exampleData = []byte(`[{
    "pciAddress": "0000:82:00.2",
    "socketpath": "/var/run/vdpa/vhost/vdpa-0"
}, {
    "pciAddress": "0000:82:00.3",
    "socketpath": "/var/run/vdpa/vhost/vdpa-1"
}, {
    "pciAddress": "0000:82:00.4",
    "socketpath": "/var/run/vdpa/vhost/vdpa-2"
}, {
    "pciAddress": "0000:82:00.5",
    "socketpath": "/var/run/vdpa/vhost/vdpa-3"
}, {
    "pciAddress": "0000:82:00.6",
    "socketpath": "/var/run/vdpa/vhost/vdpa-4"
}, {
    "pciAddress": "0000:82:00.7",
    "socketpath": "/var/run/vdpa/vhost/vdpa-5"
}, {
    "pciAddress": "0000:82:01.0",
    "socketpath": "/var/run/vdpa/vhost/vdpa-6"
}, {
    "pciAddress": "0000:82:01.1",
    "socketpath": "/var/run/vdpa/vhost/vdpa-7"
}, {
    "pciAddress": "0000:82:01.2",
    "socketpath": "/var/run/vdpa/vhost/vdpa-8"
}, {
    "pciAddress": "0000:82:01.3",
    "socketpath": "/var/run/vdpa/vhost/vdpa-9"
}, {
    "pciAddress": "0000:82:01.4",
    "socketpath": "/var/run/vdpa/vhost/vdpa-10"
}, {
    "pciAddress": "0000:82:01.5",
    "socketpath": "/var/run/vdpa/vhost/vdpa-11"
}, {
    "pciAddress": "0000:82:01.6",
    "socketpath": "/var/run/vdpa/vhost/vdpa-12"
}, {
    "pciAddress": "0000:82:01.7",
    "socketpath": "/var/run/vdpa/vhost/vdpa-13"
}, {
    "pciAddress": "0000:82:02.0",
    "socketpath": "/var/run/vdpa/vhost/vdpa-14"
}, {
    "pciAddress": "0000:82:02.1",
    "socketpath": "/var/run/vdpa/vhost/vdpa-15"
}]`)
