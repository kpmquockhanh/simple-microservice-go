package main

import (
	"context"
	"google.golang.org/grpc"
	app2 "simple-micro/core/app"
	sample_services "simple-micro/exmsg/services"
)

const (
	port = 50051
)

func main() {
	app := app2.App{
		Name: "Sample service",
		Code: "sample-service",
		Type: app2.ServiceType,
		Port: 50051,
	}
	s := grpc.NewServer()
	sample_services.RegisterSampleServer(s, &server{})
	app.NewServer(s)
}

type server struct {
	sample_services.UnimplementedSampleServer
}

func (s *server) GetNumber(ctx context.Context, req *sample_services.SampleRequest) (*sample_services.SampleResponse, error) {
	return &sample_services.SampleResponse{
		Status: true,
		Data: map[string]string{
			"number":  "123",
			"message": "Hello world",
		},
	}, nil
}
