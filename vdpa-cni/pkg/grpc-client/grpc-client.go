// SPDX-License-Identifier: Apache-2.0
// Copyright(c) 2019 Red Hat, Inc.

//
// Package vdpa manages a vDPA VF. This file contains the gRPC
// Client code. It makes gRPC calls to the vDPA DPDK gRPC Server
// to retrieve the allocated virtio socketfile based on the input
// PCI Address of the vDPA VF.
//
package vdpagrpcclient

import (
	"context"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"

	vdpagrpc "github.com/redhat-nfvpe/vdpa-deployment/grpc"
	vdpatypes "github.com/redhat-nfvpe/vdpa-deployment/pkg/types"
)

// GetSocketpath makes a gRPC call to server to find the socketpath
// associtated with the input PCI Address.
func GetSocketpath(pciAddress string) (string, error) {
	var req vdpagrpc.GetSocketpathRequest

	client, conn, err := getClient()
	if err != nil {
		log.Printf("ERROR: Error creating client - %v", err)
		return "", err
	}
	defer freeClient(conn)

	log.Printf("INFO: Getting socketpath for PCI Address %s", pciAddress)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req.PciAddress = pciAddress
	rsp, err := client.GetSocketpath(ctx, &req)
	if err != nil {
		log.Printf("ERROR: Error on gRPC call - %v: ", err)
	}
	log.Printf("INFO: Found socketpath - %s", rsp.Socketpath)

	return rsp.Socketpath, err
}

// getClient creates a new gRPC Client.
func getClient() (vdpagrpc.VdpaDpdkClient, *grpc.ClientConn, error) {
	conn, err := grpc.Dial(vdpatypes.GRPCEndpoint, grpc.WithInsecure(), grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
		grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("unix", addr, timeout)
		}),
	)

	client := vdpagrpc.NewVdpaDpdkClient(conn)

	return client, conn, err
}

// freeClient closes the gRPC Client connection.
func freeClient(conn *grpc.ClientConn) {
	conn.Close()
}

/*
// Test Data
type pciMappingTable struct {
	deviceID   string
	dpdkIndex  int
	socketpath string
}

var pciMapping = []pciMappingTable{
	pciMappingTable{
		"0000:82:00.2", 0, "/var/run/vdpa/vhost/vdpa-0",
	},
	pciMappingTable{
		"0000:82:00.3", 1, "/var/run/vdpa/vhost/vdpa-1",
	},
	pciMappingTable{
		"0000:82:00.4", 2, "/var/run/vdpa/vhost/vdpa-2",
	},
	pciMappingTable{
		"0000:82:00.5", 3, "/var/run/vdpa/vhost/vdpa-3",
	},
	pciMappingTable{
		"0000:82:00.6", 4, "/var/run/vdpa/vhost/vdpa-4",
	},
	pciMappingTable{
		"0000:82:00.7", 5, "/var/run/vdpa/vhost/vdpa-5",
	},
	pciMappingTable{
		"0000:82:01.0", 6, "/var/run/vdpa/vhost/vdpa-6",
	},
	pciMappingTable{
		"0000:82:01.1", 7, "/var/run/vdpa/vhost/vdpa-7",
	},
	pciMappingTable{
		"0000:82:01.2", 8, "/var/run/vdpa/vhost/vdpa-8",
	},
	pciMappingTable{
		"0000:82:01.3", 9, "/var/run/vdpa/vhost/vdpa-9",
	},
	pciMappingTable{
		"0000:82:01.4", 10, "/var/run/vdpa/vhost/vdpa-10",
	},
	pciMappingTable{
		"0000:82:01.5", 11, "/var/run/vdpa/vhost/vdpa-11",
	},
	pciMappingTable{
		"0000:82:01.6", 12, "/var/run/vdpa/vhost/vdpa-12",
	},
	pciMappingTable{
		"0000:82:01.7", 13, "/var/run/vdpa/vhost/vdpa-13",
	},
	pciMappingTable{
		"0000:82:02.0", 14, "/var/run/vdpa/vhost/vdpa-14",
	},
	pciMappingTable{
		"0000:82:02.1", 15, "/var/run/vdpa/vhost/vdpa-15",
	},
}

func GetSocketpath(pciAddress string) (string, error) {

	var socketpath string

	// Loop through the static PCI mapping Table and
	// find the associated socketFile.
	for _, pci := range pciMapping {
		if pciAddress == pci.deviceID {
			socketpath = pci.socketpath
			log.Printf("INFO: populateConfigData: FOUND - Mapped %s to %s", pciAddress, socketpath)
			break
		}
	}

	return socketpath, nil
}
*/
