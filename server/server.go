package server

import (
	"google.golang.org/grpc"
)

type Server struct {
	gServer *grpc.Server
	option  Option
}

var GServer *Server	// 全局服务

func Init() {
	InitConfig()

	NewRegisterContest().Register()

	NewServer(WithGRPCOpts())
}

func NewServer(opts ...Options) {
	var server Server

	for _, opt := range opts {
		opt(&server.option)
	}
	server.gServer = grpc.NewServer(server.option.gOpts...)	// 初始化grpc服务
	GServer = &server
}

func Run() {

}
