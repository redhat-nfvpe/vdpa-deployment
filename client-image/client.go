// SPDX-License-Identifier: Apache-2.0
// Copyright(c) 2019 Red Hat, Inc.

//
// Package main implements the client side for the gRPC calls.
// This code is intended as a sample implementation and used
// for testing. Not used in the vDPA deployment.
//

package main

import (
	"context"
	"flag"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"

	vdpagrpc "github.com/redhat-nfvpe/vdpa-deployment/grpc"
	vdpatypes "github.com/redhat-nfvpe/vdpa-deployment/pkg/types"
)

// getSocketpath makes a gRPC call to server to find the socketpath
// associtated with the input PCI Address.
func getSocketpath(client vdpagrpc.VdpaDpdkClient, pciAddress string) {
	var req vdpagrpc.GetSocketpathRequest

	log.Printf("INFO: Getting socketpath for PCI Address %s", pciAddress)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	req.PciAddress = pciAddress
	rsp, err := client.GetSocketpath(ctx, &req)
	if err != nil {
		log.Printf("ERROR: Unable to retrieve socketpath via gRPC - %v: ", err)
		return
	}
	log.Printf("INFO: Retrieved socketpath - %s", rsp.Socketpath)
}


func main() {
	flag.Parse()
	log.Printf("INFO: Starting vDPA-DPDK gRPC Client.")

	conn, _ := grpc.Dial(vdpatypes.GRPCEndpoint, grpc.WithInsecure(), grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
		grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("unix", addr, timeout)
		}),
	)
	defer conn.Close()

	client := vdpagrpc.NewVdpaDpdkClient(conn)

	// Call server to get associated Socketpath
	getSocketpath(client, "0000:82:00.5")
	getSocketpath(client, "0000:87:00.5")
}
