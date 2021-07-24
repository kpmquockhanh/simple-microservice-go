package app

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"
	"os/signal"
	"simple-micro/core/sd"
	"simple-micro/core/transhttp"
	"syscall"
	"time"
)

type App struct {
	Name     string
	Code     string
	BasePath string
	Port     int64
	Type     string
	ClientConnections map[string]*grpc.ClientConn
}

const APIType = "api"
const DMSType = "dms"

func (a *App) newApiServer(ctx context.Context, s transhttp.ApiServer)  {
	// Will query consul every 5 seconds.
	sd.RegisterDefault(time.Second * 5)
	a.initGrpcConnections(ctx, s.GetGrpcNames())
	s.InitGrpcClients(a.ClientConnections)
	a.initRoutes(s)
	a.HandleSigTerm(func(app *App) {
		a.CloseAllConn()
	})
}

func (a *App) newDMSServer()  {}

func (a *App) NewServer(s transhttp.ApiServer)  {
	defer a.CloseAllConn()
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	switch a.Type {
	case APIType:
		a.newApiServer(ctx, s)
	case DMSType:
		a.newDMSServer()
	}
	log.Printf("ApiServer %s started at %d", a.Type, a.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", a.Port), nil))
}

func (a *App) initRoutes(s transhttp.ApiServer)  {
	for _, route := range s.GetRoutes() {
		http.Handle(fmt.Sprintf("%s%s", a.BasePath, route.Path), route.Handler)
		log.Printf("Init route %s%s",  a.BasePath, route.Path)
	}
}

func (a *App) initGrpcConnections(ctx context.Context, clientNames []string)  {
	a.ClientConnections = make(map[string]*grpc.ClientConn)
	for _, name := range clientNames {
		address := a.getAddressByName(name)
		conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("Call grpc to %s (%s) failed: %v\n", name, address, err.Error())
		}
		a.ClientConnections[name] = conn
	}
}

func (a *App) getAddressByName(name string) string {
	return fmt.Sprintf("srv://consul/%s", name)
}

func (a *App) CloseAllConn()  {
	log.Printf("Closing all connections...")
	for _, conn := range a.ClientConnections {
		conn.Close()
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