package main

import (
	"context"
	_ "encoding/json"
	"flag"
	_ "fmt"
	_ "io"
	_ "io/ioutil"
	"log"
	_ "net"
	"time"

	"google.golang.org/grpc"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"

	_ "github.com/golang/protobuf/proto"

	vdpagrpc "github.com/redhat-nfvpe/vdpa-deployment/vdpadpdk_daemonset/grpc"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containing the CA root cert file")
	serverAddr         = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.youtube.com", "The server name use to verify the hostname returned by TLS handshake")
)


// getSocketpath makes a gRPC call to server to find the socketpath
// associtated with the input PCI Address.
func getSocketpath(client vdpagrpc.VdpaDpdkClient, pciAddress string) {
	var req vdpagrpc.GetSocketpathRequest

	log.Printf("Getting socketpath for PCI Address %s", pciAddress)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	req.PciAddress = pciAddress
	rsp, err := client.GetSocketpath(ctx, &req)
	if err != nil {
		log.Fatalf("%v.GetSocketpath(_) = _, %v: ", client, err)
	}
	log.Println(rsp.Socketpath)
}


func main() {
	flag.Parse()
	var opts []grpc.DialOption
	
	if *tls {
		if *caFile == "" {
			*caFile = testdata.Path("ca.pem")
		}
		creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := vdpagrpc.NewVdpaDpdkClient(conn)

	// Call server to get associated Socketpath
	getSocketpath(client, "0000:82:00.5")
	getSocketpath(client, "0000:87:00.5")
}
