package client

import (
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

func NewClient(opts ...Options) {
	var client Client
	client.options = &Option{
		watchInterval: 5 * time.Second,
	}
	for _, opt := range opts {
		opt(client.options)
	}
	client.connPool = make(map[string]*grpc.ClientConn)
	GClient = &client
}

func GetConn(service string) (*grpc.ClientConn, error) {
	GClient.RLock()
	if cli, ok := GClient.connPool[service]; ok {
		GClient.RUnlock()
		return cli, nil
	}
	GClient.RUnlock()

	GClient.Lock()
	defer GClient.Unlock()

	conn, err := grpc.Dial(service)
	if err != nil {
		return nil, err
	}
	GClient.connPool[service] = conn
	return conn, nil
}

func Close(service string) error {
	if conn, ok := GClient.connPool[service]; ok {
		return conn.Close()
	}
	return nil
}