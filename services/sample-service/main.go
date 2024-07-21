package main

import (
	"context"
	app2 "simple-micro/core/app"
	sample_services "simple-micro/exmsg/services"
	"simple-micro/pkg"
	"simple-micro/services/sample-service/handlers"
)

func main() {
	app := app2.App{
		Name: "Sample service",
		Code: pkg.SampleClient,
		Type: app2.ServiceType,
		Port: 50051,
	}
	s := app.NewGrpcSever()
	sample_services.RegisterSampleServer(s, &server{})
	app.NewServer(s)
}

type server struct {
	sample_services.UnimplementedSampleServer
}

func (s *server) GetNumber(ctx context.Context, req *sample_services.SampleRequest) (*sample_services.SampleResponse, error) {
	return handlers.GetNumber(ctx, req)
}
