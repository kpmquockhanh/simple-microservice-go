package main

import (
	"context"
	sample_services "simple-micro/exmsg/services"
)

type server struct {
	sample_services.UnimplementedSampleServer
}

func (s *server) GetNumber(ctx context.Context, req *sample_services.SampleRequest) (*sample_services.SampleResponse, error)  {
	return &sample_services.SampleResponse{
		Status: true,
	}, nil
}
