package sd

import (
	"strconv"
	"sync"
	"time"

	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc/resolver"
)

const (
	Scheme = "srv"
)

type ConsulResolver struct {
	lock          sync.RWMutex
	target        resolver.Target
	cc            resolver.ClientConn
	consul        *api.Client
	state         chan resolver.State
	done          chan struct{}
	watchInterval time.Duration
}

func (r *ConsulResolver) ResolveNow(resolver.ResolveNowOptions) {
	r.resolve()
}

func (r *ConsulResolver) Close() {
	close(r.done)
}

func (r *ConsulResolver) updater() {
	for {
		select {
		case state := <-r.state:
			r.cc.UpdateState(state)
		case <-r.done:
			return
		}
	}
}

func (r *ConsulResolver) watcher() {
	ticker := time.NewTicker(r.watchInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			r.resolve()
		case <-r.done:
			return
		}
	}
}

func (r *ConsulResolver) resolve() {
	r.lock.Lock()
	defer r.lock.Unlock()

	services, _, err := r.consul.Catalog().Service(r.target.Endpoint(), "", nil)
	if err != nil {
		return
	}

	addresses := make([]resolver.Address, 0, len(services))

	for _, s := range services {
		address := s.ServiceAddress
		port := s.ServicePort

		if address == "" {
			address = s.Address
		}

		addresses = append(addresses, resolver.Address{
			Addr:       address + ":" + strconv.Itoa(port),
			ServerName: r.target.Endpoint(),
		})
	}

	state := resolver.State{
		Addresses: addresses,
	}
	r.state <- state
}
