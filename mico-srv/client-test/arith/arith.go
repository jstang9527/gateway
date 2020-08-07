package main

import (
	"context"
	"fmt"

	"github.com/jstang9527/gateway/mico-srv/srv/pb"
	"google.golang.org/grpc"
)

const (
	// Address gRPC服务地址
	Address = "127.0.0.1:50055"
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
	c := pb.NewArithClient(conn)

	// 调用相加方法
	req := &pb.ArithRequest{Num1: 3, Num2: 2}
	res1, err := c.XiangJia(context.Background(), req)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("schedule add method: %v + %v = %v\n", req.Num1, req.Num2, res1.Result)
	// 调用相加方法
	res2, err := c.XiangJian(context.Background(), req)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("schedule subtract method: %v - %v = %v\n", req.Num1, req.Num2, res2.Result)
}
