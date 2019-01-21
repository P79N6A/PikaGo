package server

import (
	"fmt"
	"github.com/Carey6918/grpc/helper"
	consul "github.com/hashicorp/consul/api"
	"time"
)

/**
	使用consul进行服务发现与服务注册
	https://godoc.org/github.com/hashicorp/consul/api#pkg-index
  */

type RegisterContext struct {
	Address                        string
	ServiceName                    string
	Tags                           []string
	Port                           int
	TTL                            time.Duration
	DeregisterCriticalServiceAfter time.Duration
	Interval                       time.Duration
}

func NewRegisterContest() *RegisterContext {
	return &RegisterContext{
		Address:     "",
		ServiceName: Config.ServiceConfig.ServiceName,
		Tags:        []string{},
		Port:        Config.ServiceConfig.ServicePort,
		TTL:         2 * time.Minute,
		DeregisterCriticalServiceAfter: 1 * time.Minute,
		Interval:                       10 * time.Second,
	}
}

// 通过consul注册服务
func (r *RegisterContext) Register() error {
	config := consul.DefaultConfig()
	config.Address = r.Address
	client, err := consul.NewClient(config)
	if err != nil {
		return err
	}
	agent := client.Agent()
	localIP := helper.GetLocalIP()
	serviceID := fmt.Sprintf("%s-%d", r.ServiceName, r.Port)
	registration := &consul.AgentServiceRegistration{
		ID:      serviceID,
		Name:    r.ServiceName,
		Tags:    r.Tags,
		Port:    r.Port,
		Address: localIP,
		Check: &consul.AgentServiceCheck{ // 开启健康检查
			TTL:                            r.TTL.String(),                                          // 包体存活时间，默认为120s
			GRPC:                           fmt.Sprintf("%v:%v/%v", localIP, r.Port, r.ServiceName), //grpc 支持，执行健康检查的地址，service 会传到 Health.Check 函数中
			Interval:                       r.Interval.String(),                                     // 健康检查间隔，默认为10s
			DeregisterCriticalServiceAfter: r.DeregisterCriticalServiceAfter.String(),               // 如果检查超过这个时间，那么会自动注销这个注册
		},
	}
	return agent.ServiceRegister(registration)
}
