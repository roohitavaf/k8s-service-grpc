package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sort"
	"time"

	pb "github.com/roohitavaf/k8s-service-grpc/pkg/echo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

const (
	serviceName    = "echo-server-service"
	namespace      = "default"
	lookupInterval = 10 * time.Second
)

func connect(address string) *grpc.ClientConn {
	//with loadbalancing
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
	)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	return conn
}

func main() {
	if len(os.Args) < 2 {
		panic("provide the type of service: headless or clusterip")
	}
	serviceType := os.Args[1]
	var client pb.EchoServiceClient
	var myResolver *MyResolver
	if serviceType == "headless" {
		log.Println("Client is running to talk to a headless server...")
		myResolver = NewMyResolver(lookupInterval)
		resolver.Register(NewMyResolverBuilder(myResolver))
		fullServiceName := fmt.Sprintf("%s.%s.svc.cluster.local", serviceName, namespace)
		address := "dns://" + fullServiceName + ":50051"
		log.Println("Address: ", address)
		conn := connect(address)
		client = pb.NewEchoServiceClient(conn)

		hosts, _ := net.LookupHost(fullServiceName)
		sort.Strings(hosts)
	}

	for {
		message := "Hello, gRPC!"
		if len(os.Args) > 2 {
			message = os.Args[2]
		}
		ctx, _ := context.WithTimeout(context.Background(), time.Second)
		r, err := client.Echo(ctx, &pb.EchoRequest{Message: message})
		if err != nil {
			log.Printf("Error after calling Echo: %v", err)
			if myResolver != nil {
				myResolver.ResolveOnFailure()
			}
		} else {
			log.Printf("Response: %s", r.GetMessage())
		}
		time.Sleep(1 * time.Second)
	}

}
