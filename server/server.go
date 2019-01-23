package server

import (
	"errors"
	"fmt"
	"github.com/Carey6918/PikaRPC/helper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	gServer  *grpc.Server
	option   *Option
	listener net.Listener
}

var GServer *Server // 全局服务

func Init() {
	InitConfig()

	NewRegisterContest().Register() // 通过consul注册服务

	NewServer(WithGRPCOpts())
}

func NewServer(opts ...Options) {
	var server Server

	for _, opt := range opts {
		opt(server.option)
	}
	server.gServer = grpc.NewServer(server.option.gOpts...) // 初始化grpc服务
	GServer = &server
}

func Run() error {
	errCh := make(chan error, 1)
	go func() {
		errCh <- GServer.Serve()
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
				GServer.Stop()
				return errors.New(sig.String())
			}
		case err := <-errCh:
			return err
		}
	}
	return <-errCh
}

func (s *Server) Serve() error {
	if err := s.Listen(); err != nil {
		return err
	}

	// 注册grpc服务
	reflection.Register(s.gServer)
	if err := s.gServer.Serve(s.listener); err != nil {
		return err
	}
	return nil
}

func (s *Server) Stop() error {
	return s.listener.Close()
}

func (s *Server) Listen() error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", helper.GetLocalIP(), ServiceConf.ServicePort))
	if err != nil {
		return err
	}
	s.listener = listener
	return nil
}
