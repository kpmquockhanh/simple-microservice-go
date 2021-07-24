package main

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"net"
	"simple-micro/core/sd"
	sample_services "simple-micro/exmsg/services"
)

const (
	port = 50051
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	sample_services.RegisterSampleServer(s, &server{})

	//Use consul to register service
	service, err := sd.NewService("127.0.0.1:8500", "change-sample", port, []string{"sample_tag"})
	if err != nil {
		log.Fatalf("Failed to get new consul %v", err)
	}

	service.InitHealthCheck(s, &sd.HealthImpl{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	})

	if err := service.Register(onClose); err !=nil {
		log.Printf("Register consul failed")
		panic(err)
	}

	log.Printf("Serve on :%d", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func onClose()  {
	log.Printf("Closing...")
}