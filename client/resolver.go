package client

import (
	"github.com/Carey6918/PikaRPC/helper"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc/resolver"
	"sync"
	"time"
)

type ConsulResolver struct {
	sync.RWMutex
	target  resolver.Target
	cc      resolver.ClientConn
	client  *api.Client
	addr    chan []resolver.Address
	done    chan struct{}
	options Option
}

func (r *ConsulResolver) ResolveNow(resolver.ResolveNowOption) {
	r.resolve()
}

func (r *ConsulResolver) Close() {
	close(r.done)
}
func (r *ConsulResolver) watch() {
	ticker := time.NewTicker(r.options.watchInterval)
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

func (r *ConsulResolver) updater() {
	for {
		select {
		case addrs := <-r.addr:
			r.cc.NewAddress(addrs)
		case <-r.done:
			return
		}
	}
}

func (r *ConsulResolver) resolve() {
	r.Lock()
	defer r.Unlock()

	services, _, err := r.client.Catalog().Service(r.target.Endpoint, "", nil)
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
			Addr:       address + ":" + helper.I2S(port),
			ServerName: r.target.Endpoint,
		})
	}
	r.addr <- addresses

}
