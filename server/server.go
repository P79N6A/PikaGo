package server

import "google.golang.org/grpc"

type 适配器 struct {
	grpc服务 *grpc.Server
	选项     服务选项
}

func 生成一个适配器(选项们 ...服务选项们) *适配器 {
	var 一个适配器 适配器
	for _, 选项设置 := range 选项们 {
		选项设置(&一个适配器.选项)
	}
	一个适配器.grpc服务 = grpc.NewServer(一个适配器.选项.grpc选项...)
	return &一个适配器
}
