package server

import (
	"code.byted.org/gopkg/pkg/log"
	"fmt"
	"github.com/Carey6918/PikaRPC/helper"
	consul "github.com/hashicorp/consul/api"
	"time"
)

/**
	使用consul进行服务发现与服务注册
	https://godoc.org/github.com/hashicorp/consul/api#pkg-index
  */

const consulPort  = "8500"

type RegisterContext struct {
	ServiceName                    string
	Tags                           []string
	Port                           int
	DeregisterCriticalServiceAfter time.Duration
	Interval                       time.Duration
}

func NewRegisterContest() *RegisterContext {
	return &RegisterContext{
		ServiceName: ServiceConf.ServiceName,
		Tags:        []string{},
		Port:        helper.S2I(ServiceConf.ServicePort),
		DeregisterCriticalServiceAfter: 1 * time.Minute,
		Interval:                       10 * time.Second,
	}
}

func (r *RegisterContext) Register() error {
	config := consul.DefaultConfig()
	config.Address = helper.GetLocalAddress(consulPort)
	client, err := consul.NewClient(config)
	if err != nil {
		log.Fatal("consul new client failed, err= %v", err)
		return err
	}
	agent := client.Agent()
	localIP := helper.GetLocalIP()
	registration := &consul.AgentServiceRegistration{
		ID:      r.ServiceName,
		Name:    r.ServiceName,
		Tags:    r.Tags,
		Port:    r.Port,
		Address: localIP,
		//todo 健康检查的编码还没完成，服务在1分钟后会挂掉，所以先注释了
		Check: &consul.AgentServiceCheck{ // 开启健康检查
			GRPC:                           fmt.Sprintf("%v:%v/%v", localIP, r.Port, r.ServiceName), //grpc 支持，执行健康检查的地址，service 会传到 Health.Check 函数中
			Interval:                       r.Interval.String(),                                     // 健康检查间隔，默认为10s
			DeregisterCriticalServiceAfter: r.DeregisterCriticalServiceAfter.String(),               // 如果检查超过这个时间，那么会自动注销这个注册
		},
	}
	return agent.ServiceRegister(registration)
}
