package main

import (
	"context"
	"fmt"
	"net"

	"github.com/jstang9527/gateway/mico-srv/srv/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

// Address ...
const Address = "0.0.0.0:50055"

// ArithService ...
type ArithService struct{}

// XiangJia ...
func (a ArithService) XiangJia(ctx context.Context, req *pb.ArithRequest) (*pb.ArithResponse, error) {
	resp := new(pb.ArithResponse)
	resp.Result = req.Num1 + req.Num2
	return resp, nil
}

// XiangJian ....
func (a ArithService) XiangJian(ctx context.Context, req *pb.ArithRequest) (*pb.ArithResponse, error) {
	resp := new(pb.ArithResponse)
	resp.Result = req.Num1 - req.Num2
	return resp, nil
}

func main() {
	listen, err := net.Listen("tcp", Address)
	if err != nil {
		grpclog.Fatalf("Failed to listen: %v", err)
	}

	// 实例化grpc Server
	s := grpc.NewServer()

	// 注册HelloService
	pb.RegisterArithServer(s, ArithService{})
	// pb.RegisterHelloServer(s, ArithService)

	fmt.Println("Listen on " + Address)
	s.Serve(listen)
}
