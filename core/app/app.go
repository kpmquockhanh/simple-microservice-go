package app

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"net/http"
	"os"
	"os/signal"
	logger2 "simple-micro/core/logger"
	"simple-micro/core/sd"
	"simple-micro/core/transhttp"
	"simple-micro/core/utils"
	"syscall"
	"time"
)

type App struct {
	Name              string
	Code              string
	BasePath          string
	Port              int64
	Type              string
	ClientConnections map[string]*grpc.ClientConn
	logger            *zap.SugaredLogger
}

const APIType = "api"
const ServiceType = "service"

func loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	logger := logger2.NewLogger()
	logger.Infof("Received request: %v with req: %s", info.FullMethod, utils.ToJsonString(req))
	resp, err := handler(ctx, req)
	return resp, err
}

func (a *App) NewGrpcSever() *grpc.Server {
	return grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor),
	)
}

func (a *App) newApiServer(ctx context.Context, s transhttp.ApiServer) {
	// Will query consul every 5 seconds.
	sd.RegisterDefault(time.Second * 5)
	a.initGrpcConnections(ctx, s.GetGrpcNames())
	s.InitGrpcClients(a.ClientConnections)
	a.initRoutes(s)
	a.HandleSigTerm(func(app *App) {
		a.logger.Info("Closing...")
		a.CloseAllConn()
	})
	logger := a.logger

	logger.Infof("Api server %s started at %d", a.Type, a.Port)
	logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", a.Port), nil))
}

func (a *App) newServiceServer(s *grpc.Server) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.Port))
	if err != nil {
		a.logger.Fatalf("failed to listen: %v", err)
	}
	//Use consul to register service
	service, err := sd.NewService(fmt.Sprintf("consul:8500"), a.Code, int(a.Port), []string{"sample_tag"})
	if err != nil {
		a.logger.Fatalf("Failed to get new consul %v", err)
	}

	service.InitHealthCheck(s, &sd.HealthImpl{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	})

	if err := service.Register(a.onClose); err != nil {
		a.logger.Infof("Register consul failed")
		panic(err)
	}

	a.logger.Infof("Service %s started at %d", a.Code, a.Port)
	if err := s.Serve(lis); err != nil {
		a.logger.Fatalf("failed to serve: %v", err)
	}
}

func (a *App) NewServer(s interface{}) {
	a.logger = logger2.NewLogger()
	defer a.CloseAllConn()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	switch a.Type {
	case APIType:
		a.newApiServer(ctx, s.(transhttp.ApiServer))
	case ServiceType:
		a.newServiceServer(s.(*grpc.Server))
	default:
		panic("invalid server type")
	}
}

func (a *App) initRoutes(s transhttp.ApiServer) {
	for _, route := range s.GetRoutes() {
		a.logger.Infof("Init route %s%s", a.BasePath, route.Path)
		http.Handle(fmt.Sprintf("%s%s", a.BasePath, route.Path), route.Handler)
	}
}

func (a *App) initGrpcConnections(ctx context.Context, clientNames []string) {
	a.ClientConnections = make(map[string]*grpc.ClientConn)
	for _, name := range clientNames {
		address := a.getAddressByName(name)
		a.logger.Infof("Dialing grpc to %s", address)
		conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			a.logger.Fatalf("Call grpc to %s (%s) failed: %v\n", name, address, err.Error())
		}
		a.ClientConnections[name] = conn
	}
}

func (a *App) getAddressByName(name string) string {
	return fmt.Sprintf("srv://consul/%s", name)
}

func (a *App) CloseAllConn() {
	a.logger.Infof("Closing all connections...")
	for _, conn := range a.ClientConnections {
		_ = conn.Close()
	}
}

func (a *App) HandleSigTerm(onClose func(app *App)) {
	done := make(chan os.Signal)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-done
		onClose(a)
		os.Exit(0)
	}()
}

func (a *App) onClose() {
	a.logger.Infow("Closing...")
}
