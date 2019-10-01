package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	_ "io"
	_ "io/ioutil"
	"log"
	"net"

	"google.golang.org/grpc"

	_ "google.golang.org/grpc/credentials"
	_ "google.golang.org/grpc/testdata"

	_ "github.com/golang/protobuf/proto"

	vdpagrpc "github.com/redhat-nfvpe/vdpa-deployment/grpc"
)

var (
	//tls        = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	//certFile   = flag.String("cert_file", "", "The TLS cert file")
	//keyFile    = flag.String("key_file", "", "The TLS key file")
	//jsonDBFile = flag.String("json_db_file", "", "A json file containing a list of features")
	port       = flag.Int("port", 10000, "The server port")
)

type socketMapping struct {
	PciAddress string  `json:"pciAddress"`
	DpdkIndex  int     `json:"dpdkIndex"`
}

type vdpaDpdkServer struct {
	mappingTable []*socketMapping
	baseSocketpath string
}

// GetSocketpath takes a PCI Address of a vDPA VF and returns the
// associated Socketpath
func (s *vdpaDpdkServer) GetSocketpath(ctx context.Context, req *vdpagrpc.GetSocketpathRequest) (*vdpagrpc.GetSocketpathResponse, error) {
	var rsp vdpagrpc.GetSocketpathResponse

	log.Printf("Received request for PCI Address %s", req.PciAddress)
	for _, vdpaIface := range s.mappingTable {
		if vdpaIface.PciAddress == req.PciAddress {
			rsp.Socketpath = fmt.Sprintf("%s%d", s.baseSocketpath, vdpaIface.DpdkIndex)
			log.Printf("Found File: %s", rsp.Socketpath)
			break
		}
	}
	if rsp.Socketpath == "" {
		log.Printf("No File Found!")
	}
	
	return &rsp, nil
}

// scanInterfaces scans the system for available vDPA Interfaces (PCI Addresses).
func (s *vdpaDpdkServer) scanInterfaces() {
	if err := json.Unmarshal(exampleData, &s.mappingTable); err != nil {
		log.Fatalf("Failed to load default mappingTable: %v", err)
	}
	s.baseSocketpath = "/var/run/vdpa/vdpa-"

	log.Printf("Loaded PCI Addresses: %s", s.baseSocketpath)
	for i, vdpaIface := range s.mappingTable {
		log.Printf("  %d - %s - Index %d", i, vdpaIface.PciAddress, vdpaIface.DpdkIndex)
	}
}

// newServer creates a new server instance with initialized data.
func newServer() *vdpaDpdkServer {
	s := &vdpaDpdkServer{}
	s.scanInterfaces()
	return s
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	/*
	if *tls {
		if *certFile == "" {
			*certFile = testdata.Path("server1.pem")
		}
		if *keyFile == "" {
			*keyFile = testdata.Path("server1.key")
		}
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			log.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	*/
	grpcServer := grpc.NewServer(opts...)
	vdpagrpc.RegisterVdpaDpdkServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}

// exampleData is a copy of vdpa_db.json. It's to avoid
// specifying file path with `go run`.
var exampleData = []byte(`[{
    "pciAddress": "0000:82:00.2",
    "dpdkIndex": 0
}, {
    "pciAddress": "0000:82:00.3",
    "dpdkIndex": 1
}, {
    "pciAddress": "0000:82:00.4",
    "dpdkIndex": 2
}, {
    "pciAddress": "0000:82:00.5",
    "dpdkIndex": 3
}, {
    "pciAddress": "0000:82:00.6",
    "dpdkIndex": 4
}, {
    "pciAddress": "0000:82:00.7",
    "dpdkIndex": 5
}, {
    "pciAddress": "0000:82:01.0",
    "dpdkIndex": 6
}, {
    "pciAddress": "0000:82:01.1",
    "dpdkIndex": 7
}, {
    "pciAddress": "0000:82:01.2",
    "dpdkIndex": 8
}, {
    "pciAddress": "0000:82:01.3",
    "dpdkIndex": 9
}, {
    "pciAddress": "0000:82:01.4",
    "dpdkIndex": 10
}, {
    "pciAddress": "0000:82:01.5",
    "dpdkIndex": 11
}, {
    "pciAddress": "0000:82:01.6",
    "dpdkIndex": 12
}, {
    "pciAddress": "0000:82:01.7",
    "dpdkIndex": 13
}, {
    "pciAddress": "0000:82:02.0",
    "dpdkIndex": 14
}, {
    "pciAddress": "0000:82:02.1",
    "dpdkIndex": 15
}]`)
