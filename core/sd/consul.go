package sd

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// ConsulRegister consul service register
type ConsulService struct {
	Client                         *api.Client
	Address                        string
	Name                           string
	Tag                            []string
	Port                           int
	DeregisterCriticalServiceAfter time.Duration
	Interval                       time.Duration
	ID                             string
}

func NewService(addr string, name string, port int, tags []string) (*ConsulService, error) {
	config := api.DefaultConfig()
	config.Address = addr
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &ConsulService{
		Client:                         client,
		Address:                        addr,
		Name:                           name,
		Port:                           port,
		Tag:                            tags,
		DeregisterCriticalServiceAfter: time.Duration(10) * time.Second,
		Interval:                       time.Duration(10) * time.Second,
	}, nil
}

func (r *ConsulService) InitHealthCheck(s *grpc.Server, healthSrv *HealthImpl) {
	grpc_health_v1.RegisterHealthServer(s, healthSrv)
}

func (r *ConsulService) GetAddress(name string) string {
	return name
}

// Register register service
func (r *ConsulService) Register(onClose func()) error {
	agent := r.Client.Agent()
	IP := LocalIP()
	ID := fmt.Sprintf("go.sample.grpc.%v-%v-%v", r.Name, IP, r.Port)
	reg := &api.AgentServiceRegistration{
		ID:      ID, // Name of the service node
		Name:    r.Name,
		Tags:    r.Tag,
		Port:    r.Port,
		Address: IP,
		Check: &api.AgentServiceCheck{
			Interval:                       r.Interval.String(),
			GRPC:                           fmt.Sprintf("%v:%v/%v", IP, r.Port, r.Name),
			DeregisterCriticalServiceAfter: r.DeregisterCriticalServiceAfter.String(),
		},
	}
	r.handleSigterm(onClose)

	if err := agent.ServiceRegister(reg); err != nil {
		return err
	}

	r.ID = ID
	return nil
}

// DeRegister register service
func (r *ConsulService) DeRegister() error {
	agent := r.Client.Agent()
	if err := agent.ServiceDeregister(r.ID); err != nil {
		return err
	}

	return nil
}

func (r *ConsulService) DeRegisterByID(ID string) error {
	agent := r.Client.Agent()
	if err := agent.ServiceDeregister(ID); err != nil {
		return err
	}

	return nil
}

func (r *ConsulService) handleSigterm(onClose func()) {
	done := make(chan os.Signal)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-done
		// DeRegister service
		err := r.DeRegister()
		if err != nil {
			log.Printf("Deregister error service ID: %v", r.ID)
			log.Fatal(err)
		}
		log.Printf("DeRegistered service %v", r.ID)
		onClose()
		os.Exit(0)
	}()
}

func LocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
