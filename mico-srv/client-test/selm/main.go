package main

import (
	"context"
	"fmt"

	pb "github.com/jstang9527/gateway/mico-srv/srv/pb"
	"google.golang.org/grpc"
)

const (
	// Address gRPC服务地址
	Address = "127.0.0.1:50052"
)

func main() {
	// 连接
	conn, err := grpc.Dial(Address, grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	// 初始化客户端
	c := pb.NewSeleniumClient(conn)

	// 调用方法
	req := &pb.SeleniumRequest{Url: "xxx", SearchTimeout: 10} // &pb.HelloRequest{Name: "gRPC"}
	res, err := c.RunTest(context.Background(), req)

	if err != nil {
		fmt.Println(err)
	}
	if res.Message == "" {
		fmt.Println("success.")
	} else {
		fmt.Println("failed. info: ", err)
	}
}
