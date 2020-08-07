package main

import (
	"context"
	"fmt"

	"github.com/jstang9527/gateway/mico-srv/srv/pb"
	"google.golang.org/grpc"
)

const (
	// Address gRPC服务地址
	Address = "127.0.0.1:50061"
)

func main() {
	// 连接
	conn, err := grpc.Dial(Address, grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	// 初始化客户端
	// c := pb.NewHelloClient(conn)
	c := pb.NewVMClient(conn)
	// ------------------------------获取domain列表

	// ------------------------------调用开机方法
	req := &pb.VMRequest{Domain: "centos775"}
	resp, err := c.Start(context.Background(), req)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("schedule shutdown method: %v-%v\n", resp.Status, resp.Errmsg)
	// ------------------------------调用关机方法
	// resp, err := c.ShutDown(context.Background(), req)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Printf("schedule shutdown method: %v-%v\n", resp.Status, resp.Errmsg)
}
