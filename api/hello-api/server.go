package main

import (
	"google.golang.org/grpc"
	"log"
	"simple-micro/api/hello-api/handler"
	"simple-micro/core/transhttp"
	sample_services "simple-micro/exmsg/services"
	"simple-micro/pkg"
)

type server struct {
	SampleClient sample_services.SampleClient
}

func (s *server) GetGrpcNames() []string {
	return []string{
		pkg.SampleClient,
	}
}

func (s *server) GetRoutes() transhttp.Routes {
	return transhttp.Routes{
		transhttp.Route{
			Name:    "Hello",
			Path:    "/demo",
			Handler: &handler.HelloHandler{
				SampleClient: s.SampleClient,
			},
		},
		transhttp.Route{
			Name:    "Hello2",
			Path:    "/test",
			Handler: &handler.HelloHandler{},
		},
	}
}

func (s *server) InitGrpcClients(clientConn map[string]*grpc.ClientConn)  {
	conn := clientConn[pkg.SampleClient]
	s.SampleClient = sample_services.NewSampleClient(conn)
	log.Printf("Init grpc done")
}