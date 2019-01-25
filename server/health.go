package server

import (
	"code.byted.org/gopkg/pkg/log"
	"context"
	"github.com/Carey6918/PikaRPC/client"
	health "google.golang.org/grpc/health/grpc_health_v1"
	"time"
)

// gRPC健康检查，实现了grpc_health_v1.HealthServer接口
type HealthServerImpl struct{}

func (s *HealthServerImpl) Check(ctx context.Context, req *health.HealthCheckRequest) (*health.HealthCheckResponse, error) {
	client.Init(client.WithWatchInterval(10 * time.Second))
	_, err := client.GetConn(req.GetService())
	defer client.Close(req.GetService())
	if err != nil {
		log.Errorf("health check to %v failed, err= %v", req.GetService(), err)
		return &health.HealthCheckResponse{
			Status: health.HealthCheckResponse_NOT_SERVING,
		}, nil
	}
	log.Infof("health check to %v success", req.GetService())
	return &health.HealthCheckResponse{
		Status: health.HealthCheckResponse_SERVING,
	}, nil
}

func (s *HealthServerImpl) Watch(req *health.HealthCheckRequest, server health.Health_WatchServer) error {
	return nil
}
