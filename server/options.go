package server

import "google.golang.org/grpc"

type Option struct {
	gOpts []grpc.ServerOption
}

type Options func(o *Option)

func WithGRPCOpts(gopts ...grpc.ServerOption) Options {
	return func(o *Option) {
		o.gOpts = gopts
	}
}
