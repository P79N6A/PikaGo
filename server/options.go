package server

import "google.golang.org/grpc"

type 服务选项 struct {
	grpc选项 []grpc.ServerOption
}

type 服务选项们 func(o *服务选项)
