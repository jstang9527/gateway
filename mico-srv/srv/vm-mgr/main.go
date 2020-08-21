package main

import (
	"context"
	"fmt"
	"net"
	"time"

	qemu "github.com/jstang9527/gateway/mico-srv/modules"
	"github.com/jstang9527/gateway/mico-srv/srv/pb"
	"github.com/jstang9527/opentest/ops"
	"google.golang.org/grpc"
)

// Address ...
const Address = "0.0.0.0:50061"

// VMService ...
type VMService struct{}

// Start ...
func (a VMService) Start(ctx context.Context, req *pb.VMRequest) (*pb.VMResponse, error) {
	resp := new(pb.VMResponse)
	resp.Status = true

	//获取主机domain
	vmObj := qemu.NewVM(req.Domain)
	if err := vmObj.Start(); err != nil { //这个是还原
		resp.Status = false
		resp.Errmsg = fmt.Sprint(err)
	}
	return resp, nil
}

// ShutDown 不管它存在与否,直接关
func (a VMService) ShutDown(ctx context.Context, req *pb.VMRequest) (*pb.VMResponse, error) {
	resp := new(pb.VMResponse)
	resp.Status = true

	// 获取主机domain
	vmObj := qemu.NewVM(req.Domain)
	if err := vmObj.ShutDown(); err != nil {
		resp.Status = false
		resp.Errmsg = fmt.Sprint(err)
	}
	return resp, nil
}

// Recovery ...
func (a VMService) Recovery(ctx context.Context, req *pb.VMRequest) (*pb.VMResponse, error) {
	resp := new(pb.VMResponse)
	resp.Status = true
	resp.Running = true
	// 获取主机domain
	vmObj := qemu.NewVM(req.Domain)
	if err := vmObj.Recover(); err != nil {
		resp.Status = false
		resp.Errmsg = fmt.Sprint(err)
		resp.Running = false
	}
	return resp, nil
}

// Running ...
func (a VMService) Running(ctx context.Context, req *pb.VMRequest) (*pb.VMResponse, error) {
	resp := new(pb.VMResponse)

	// 获取主机domain
	vmObj := qemu.NewVM(req.Domain)
	status, err := vmObj.IsRunning()
	if err != nil {
		resp.Errmsg = fmt.Sprint(err)
		return resp, err
	}
	resp.Status = status
	resp.Running = status
	return resp, nil
}

func main() {
	// 连接本机的Qemu
	if err := qemu.Init("unix", "/var/run/libvirt/libvirt-sock", time.Second*2); err != nil {
		ops.Console("failed to connect to hypervisor: %v", err)
		return
	}
	ops.Console("success connect to hypervisor.")
	listen, err := net.Listen("tcp", Address)
	if err != nil {
		ops.Console("Failed to listen: %v", err)
		return
	}
	defer listen.Close()

	// 实例化grpc Server
	s := grpc.NewServer()

	// 注册Service
	pb.RegisterVMServer(s, VMService{})

	ops.Console("Listen on " + Address)
	s.Serve(listen)
}
