package main

import (
	"google.golang.org/grpc"
	"simple-micro/api/hello-api/handler"
	app2 "simple-micro/core/app"
	logger2 "simple-micro/core/logger"
	"simple-micro/core/transhttp"
	sample_services "simple-micro/exmsg/services"
	"simple-micro/pkg"
)

func main() {
	app := app2.App{
		Name:     "Hello",
		Code:     "hello",
		Type:     app2.APIType,
		BasePath: "/hello",
		Port:     50052,
	}

	app.NewServer(&server{})
}

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
			Name: "Hello",
			Path: "/demo",
			Handler: &handler.HelloHandler{
				SampleClient: s.SampleClient,
			},
		},
	}
}

func (s *server) InitGrpcClients(clientConn map[string]*grpc.ClientConn) {
	conn := clientConn[pkg.SampleClient]
	s.SampleClient = sample_services.NewSampleClient(conn)
	logger := logger2.NewLogger()
	logger.Info("Init grpc done")
}
