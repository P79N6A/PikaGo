package server

import "google.golang.org/grpc"

type Option struct {
	gOpts []grpc.ServerOption
}

type Options func(o *Option)

func WithGRPCOpts(gOpts ...grpc.ServerOption) Options {
	return func(o *Option) {
		o.gOpts = gOpts
	}
}
