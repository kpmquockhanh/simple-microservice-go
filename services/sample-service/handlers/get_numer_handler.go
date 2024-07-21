package handlers

import (
	"context"
	"github.com/spf13/cast"
	logger2 "simple-micro/core/logger"
	sample_services "simple-micro/exmsg/services"
)

func GetNumber(ctx context.Context, req *sample_services.SampleRequest) (*sample_services.SampleResponse, error) {
	logger := logger2.NewLogger()
	logger.Infof("Received %d", req.GetId())
	return &sample_services.SampleResponse{
		Status: true,
		Data: map[string]string{
			"number": cast.ToString(req.GetId()),
		},
	}, nil
}
