package client

import (
	consul "github.com/hashicorp/consul/api"
	"google.golang.org/grpc/resolver"
)

type ConsulBuilder struct {
	scheme string
}

func NewBuilder(scheme string) resolver.Builder {
	return &ConsulBuilder{
		scheme: scheme,
	}
}
func (b *ConsulBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {

	client, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		return nil, err
	}
	r := &ConsulResolver{
		target: target,
		cc:     cc,
		client: client,
		addr:   make(chan []resolver.Address, 1),
		done:   make(chan struct{}, 1),
	}
	go r.updater()
	go r.watch()
	r.resolve()

	return r, nil
}

func (b *ConsulBuilder) Scheme() string {
	return b.scheme
}
