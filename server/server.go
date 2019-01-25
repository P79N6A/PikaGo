package server

import (
	"code.byted.org/gopkg/pkg/log"
	"errors"
	"fmt"
	"github.com/Carey6918/PikaRPC/helper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	gServer  *grpc.Server
	option   *Option
	listener net.Listener
}

var GServer *Server // 全局服务

func Init() {
	InitConfig()

	// 通过consul注册服务
	if err := NewRegisterContest().Register(); err != nil {
		log.Fatalf("consul register failed, err= %v", err)
	}


	NewServer(WithGRPCOpts(grpc.ConnectionTimeout(500 * time.Millisecond)))
	grpc_health_v1.RegisterHealthServer(GetGRPCServer(), &HealthServerImpl{})
}

func NewServer(opts ...Options) {
	var server Server
	server.option = new(Option)
	for _, opt := range opts {
		opt(server.option)
	}
	// 初始化gRPC服务
	server.gServer = grpc.NewServer(server.option.gOpts...)
	GServer = &server
}

func Run() error {
	errCh := make(chan error, 1)
	go func() {
		errCh <- GServer.serve()
	}()
	return waitSignal(errCh)
}

func waitSignal(errCh chan error) error {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)

	for {
		select {
		case sig := <-signals:
			switch sig {
			// exit forcely
			case syscall.SIGTERM: // 结束程序(可以被捕获、阻塞或忽略)
				return errors.New(sig.String())
			case syscall.SIGHUP, syscall.SIGINT: // 终端连接断开/用户发送(ctrl+c)结束
				GServer.stop()
				return errors.New(sig.String())
			}
		case err := <-errCh:
			return err
		}
	}
	return <-errCh
}

func (s *Server) serve() error {
	if err := s.listen(); err != nil {
		return err
	}

	// 注册gRPC服务
	reflection.Register(s.gServer)
	if err := s.gServer.Serve(s.listener); err != nil {
		log.Errorf("grpc serve failed, err= %v", err)
		return err
	}
	return nil
}

func (s *Server) stop() error {
	return s.listener.Close()
}

func (s *Server) listen() error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", helper.GetLocalIP(), ServiceConf.ServicePort))
	if err != nil {
		log.Errorf("listen tcp failed, err= %v", err)
		return err
	}
	s.listener = listener
	return nil
}

func GetGRPCServer() *grpc.Server {
	return GServer.gServer
}
