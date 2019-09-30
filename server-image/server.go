package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	_ "io"
	"io/ioutil"
	"log"
	"net"

	"google.golang.org/grpc"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"

	_ "github.com/golang/protobuf/proto"

	vdpagrpc "github.com/redhat-nfvpe/vdpa-deployment/grpc"
	vdpatypes "github.com/redhat-nfvpe/vdpa-deployment/pkg/types"
)

var (
	tls        = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile   = flag.String("cert_file", "", "The TLS cert file")
	keyFile    = flag.String("key_file", "", "The TLS key file")
	jsonDBFile = flag.String("json_db_file", "", "A json file containing a list of vDPA VF Interfaces")
	port       = flag.Int("port", 10000, "The server port")
)

type vdpaDpdkServer struct {
	mappingTable []*vdpatypes.VFPCIMapping
}

// GetSocketpath takes a PCI Address of a vDPA VF and returns the
// associated Socketpath
func (s *vdpaDpdkServer) GetSocketpath(ctx context.Context, req *vdpagrpc.GetSocketpathRequest) (*vdpagrpc.GetSocketpathResponse, error) {
	var rsp vdpagrpc.GetSocketpathResponse

	log.Printf("Received request for PCI Address %s", req.PciAddress)
	for _, vdpaIface := range s.mappingTable {
		if vdpaIface.PciAddress == req.PciAddress {
			rsp.Socketpath = vdpaIface.Socketpath
			log.Printf("Found File: %s", rsp.Socketpath)
			break
		}
	}
	if rsp.Socketpath == "" {
		log.Printf("No File Found!")
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
			log.Fatalf("Failed to load default features: %v", err)
		}
	} else {
		data = exampleData
	}
	if err := json.Unmarshal(data, &s.mappingTable); err != nil {
		log.Fatalf("Failed to load default features: %v", err)
	}
}

// newServer creates a new server instance with initialized data.
func newServer() *vdpaDpdkServer {
	s := &vdpaDpdkServer{}
	s.loadInterfaces(*jsonDBFile)
	return s
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
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
	grpcServer := grpc.NewServer(opts...)
	vdpagrpc.RegisterVdpaDpdkServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}


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
