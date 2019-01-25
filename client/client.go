package client

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"sync"
	"time"
)

func init() {
	resolver.Register(NewBuilder("test")) // consul lb
}

type Client struct {
	sync.RWMutex
	connPool map[string]*grpc.ClientConn
	options  *Option
}

var GClient *Client

func Init(opts ...Options) {
	var client Client
	client.options = &Option{
		watchInterval: 20 * time.Second,
	}
	for _, opt := range opts {
		opt(client.options)
	}
	client.connPool = make(map[string]*grpc.ClientConn)
	GClient = &client
}

func GetConn(serviceName string) (*grpc.ClientConn, error) {
	GClient.RLock()
	if cli, ok := GClient.connPool[serviceName]; ok {
		GClient.RUnlock()
		return cli, nil
	}
	GClient.RUnlock()

	// 通过consul服务发现
	service, err := discovery(serviceName)
	if err != nil {
		return nil, err
	}
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", service.Address, service.Port), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	GClient.Lock()
	defer GClient.Unlock()
	GClient.connPool[serviceName] = conn
	return conn, nil
}

func Close(service string) error {
	if conn, ok := GClient.connPool[service]; ok {
		return conn.Close()
	}
	return nil
}
