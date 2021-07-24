package transhttp

import "google.golang.org/grpc"

type ApiServer interface {
	GetRoutes() Routes
	GetGrpcNames() []string
	InitGrpcClients(clientConn map[string]*grpc.ClientConn)
}

type DmsServer interface {
	RegisterServer()
}