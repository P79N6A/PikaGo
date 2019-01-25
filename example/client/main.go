package main

import (
	"context"
	"github.com/Carey6918/PikaRPC/client"
	"github.com/Carey6918/PikaRPC/example/proto"
	"log"
	"time"
)

const ServiceName = "carey.is.genius"

func main() {
	client.Init(client.WithWatchInterval(10 * time.Second))
	conn, err := client.GetConn(ServiceName)
	defer client.Close(ServiceName)
	if err != nil {
		log.Fatal(err)
	}

	cli := add.NewAddServiceClient(conn)
	req := &add.AddRequest{
		A: 1,
		B: 1,
	}
	resp, err := cli.Add(context.Background(), req)
	if err != nil {
		log.Fatalf("add failed, err= %v", err)
	}
	log.Printf("resp= %v", resp)
}
