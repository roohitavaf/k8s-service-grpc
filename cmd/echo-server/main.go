package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"sync"
	"time"

	pb "github.com/roohitavaf/k8s-service-grpc/pkg/echo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type echoServer struct {
	pb.UnimplementedEchoServiceServer
	ready        bool
	mu           sync.RWMutex
	unreadyTimer *time.Timer
	ip           string
}

func (s *echoServer) Echo(ctx context.Context, req *pb.EchoRequest) (*pb.EchoResponse, error) {
	s.mu.RLock()
	if !s.ready {
		s.mu.RUnlock()
		return nil, grpc.Errorf(codes.Aborted, "Server is not ready")
	}
	s.mu.RUnlock()
	log.Printf("Received(%v): %v", s.ip, req.GetMessage())

	// Randomly decide if the server should become unready
	if rand.Float32() < 0.2 {
		s.mu.Lock()
		s.ready = false
		s.mu.Unlock()
		log.Println("Server is now unready.")
		// Start a timer to become ready again after 5 seconds
		s.startUnreadyTimer()
	}

	return &pb.EchoResponse{Message: "Server_IP(" + s.ip + "): " + req.GetMessage()}, nil
}

func (s *echoServer) readinessHandler(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	if s.ready {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Ready"))
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("Not Ready"))
	}
	s.mu.RUnlock()
}

func (s *echoServer) startUnreadyTimer() {
	if s.unreadyTimer != nil {
		s.unreadyTimer.Stop()
	}

	randomDuration := time.Duration(rand.Intn(20)+1) * time.Second

	s.unreadyTimer = time.AfterFunc(randomDuration, func() {
		log.Printf("Going to make the server ready again...")
		s.mu.Lock()
		s.ready = true
		s.mu.Unlock()
		log.Println("Server is now ready again.")
	})
	log.Printf("Server will be ready again in %v.", randomDuration)
}

func getLocalIP() (string, error) {
	var ipAddress string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			ipAddress = ipnet.IP.String()
			break
		}
	}
	if ipAddress == "" {
		return "", fmt.Errorf("no valid IP address found")
	}
	return ipAddress, nil
}

func main() {
	rand.Seed(time.Now().UnixNano())

	ipAddress, err := getLocalIP()
	if err != nil {
		log.Fatalf("Failed to get IP address: %v", err)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	echoSvc := &echoServer{
		ready: true,
		ip:    ipAddress,
	}
	pb.RegisterEchoServiceServer(s, echoSvc)

	// Start readiness HTTP server
	go func() {
		http.HandleFunc("/readiness", echoSvc.readinessHandler)
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("Failed to start readiness server: %v", err)
		}
	}()

	log.Println("Server is running on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
