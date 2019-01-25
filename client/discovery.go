package client

import (
	"code.byted.org/golf/consul"
	"errors"
	"github.com/Carey6918/PikaRPC/helper"
)

const consulPort = "8500"

func discovery(serviceName string) (*consul.AgentService, error) {
	config := consul.DefaultConfig()
	config.Address = helper.GetLocalAddress(consulPort)
	client, err := consul.NewClient(config)
	if err != nil {
		return nil, err
	}
	services, err := client.Agent().Services()
	service, ok := services[serviceName]
	if !ok {
		return nil, errors.New("cannot find service")
	}
	return service, nil
}
