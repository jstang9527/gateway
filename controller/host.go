package controller

import (
	"context"
	"fmt"

	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"github.com/jstang9527/gateway/dao"
	"github.com/jstang9527/gateway/dto"
	"github.com/jstang9527/gateway/mico-srv/srv/pb"
	"github.com/jstang9527/gateway/middleware"
	"google.golang.org/grpc"
)

// HostController 主机控制器结果体
type HostController struct{}

// HostRegister 主机控制器
func HostRegister(group *gin.RouterGroup) {
	HostCtl := &HostController{}
	group.GET("/host_list", HostCtl.HostList)              //记录列表
	group.POST("/host", nil)                               //新增记录
	group.PUT("/host", nil)                                //修改记录
	group.DELETE("/host", nil)                             //移除记录
	group.GET("/control/status", HostCtl.GetDoaminStatus)  //获取域状态
	group.PUT("/control/start", HostCtl.StartDoamin)       //启动域
	group.PUT("/control/shutdown", HostCtl.ShutdownDoamin) //关闭域
	group.PUT("/control/recover", HostCtl.RecoverDoamin)   //还原域
}

// HostList godoc
// @Summary 主机列表
// @Description 主机列表
// @Tags 主机管理
// @ID /host/host_list
// @Accept json
// @Produce json
// @Param info query string false "关键词"
// @Param page_size query int true "每页个数"
// @Param page_no query int true "当前页数"
// @Success 200 {object} middleware.Response{data=dto.HostListOutput} "success"
// @Router /host/host_list [get]
func (s *HostController) HostList(c *gin.Context) {
	inputParams := &dto.HostListInput{}
	if err := inputParams.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	//从db中分页读取基本信息
	hostInfo := &dao.HostInfo{}
	hostList, total, err := hostInfo.PageList(c, tx, inputParams)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	//格式化输出信息
	outList := []dto.HostListItemOutput{} //这个结构体是面向前端接口的
	for _, item := range hostList {
		outItem := dto.HostListItemOutput{
			ID:         item.ID,
			Domain:     item.Domain,
			DomainIP:   item.DomainIP,
			DomainOS:   item.DomainOS,
			DomainType: item.DomainType,
			DomainDesc: item.DomainDesc,
			HostIP:     item.HostIP,
			HostDesc:   item.HostDesc,
		}
		outList = append(outList, outItem)
	}
	out := &dto.HostListOutput{Total: total, List: outList}
	middleware.ResponseSuccess(c, out)
}

// StartDoamin godoc
// @Summary 启动域
// @Description 启动域
// @Tags 主机管理
// @ID /host/control/start
// @Accept json
// @Produce json
// @Param domain query string true "域名"
// @Param host_ip query string true "宿主机ip"
// @Success 200 {object} middleware.Response{data=dto.HostStatusOutput} "success"
// @Router /host/control/start [put]
func (s *HostController) StartDoamin(c *gin.Context) {
	inputParams := &dto.HostInfoInput{}
	if err := inputParams.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	// + GRPC 调用 ----------------------
	// 连接
	srvAddr := fmt.Sprintf("%s:%v", inputParams.HostIP, 50061)
	conn, err := grpc.Dial(srvAddr, grpc.WithInsecure())
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	defer conn.Close()

	// 初始化客户端
	instance := pb.NewVMClient(conn)
	// c := pb.NewVMClient(conn)
	// 调用开机方法
	req := &pb.VMRequest{Domain: inputParams.Domain}
	resp, err := instance.Start(context.Background(), req)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	// ------------------------------完
	out := &dto.HostStatusOutput{Domain: inputParams.Domain, HostIP: inputParams.HostIP, Status: resp.Status, Errmsg: resp.Errmsg}
	middleware.ResponseSuccess(c, out)
}

// ShutdownDoamin godoc
// @Summary 关闭域
// @Description 关闭域
// @Tags 主机管理
// @ID /host/control/shutdown
// @Accept json
// @Produce json
// @Param domain query string true "域名"
// @Param host_ip query string true "宿主机ip"
// @Success 200 {object} middleware.Response{data=dto.HostStatusOutput} "success"
// @Router /host/control/shutdown [put]
func (s *HostController) ShutdownDoamin(c *gin.Context) {
	inputParams := &dto.HostInfoInput{}
	if err := inputParams.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	// + GRPC 调用 ----------------------
	// 连接
	srvAddr := fmt.Sprintf("%s:%v", inputParams.HostIP, 50061)
	conn, err := grpc.Dial(srvAddr, grpc.WithInsecure())
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	defer conn.Close()

	// 初始化客户端
	instance := pb.NewVMClient(conn)
	// c := pb.NewVMClient(conn)
	// 调用关机方法
	req := &pb.VMRequest{Domain: inputParams.Domain}
	resp, err := instance.ShutDown(context.Background(), req)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	// ------------------------------完
	out := &dto.HostStatusOutput{Domain: inputParams.Domain, HostIP: inputParams.HostIP, Status: resp.Status, Errmsg: resp.Errmsg}
	middleware.ResponseSuccess(c, out)
}

// GetDoaminStatus godoc
// @Summary 域状态
// @Description 域状态
// @Tags 主机管理
// @ID /host/control/status
// @Accept json
// @Produce json
// @Param domain query string true "域名"
// @Param host_ip query string true "宿主机ip"
// @Success 200 {object} middleware.Response{data=dto.HostStatusOutput} "success"
// @Router /host/control/status [get]
func (s *HostController) GetDoaminStatus(c *gin.Context) {
	inputParams := &dto.HostInfoInput{}
	if err := inputParams.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	// + GRPC 调用 ----------------------
	// 连接
	srvAddr := fmt.Sprintf("%s:%v", inputParams.HostIP, 50061)
	conn, err := grpc.Dial(srvAddr, grpc.WithInsecure())
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	defer conn.Close()

	// 初始化客户端
	instance := pb.NewVMClient(conn)
	// c := pb.NewVMClient(conn)
	// 调用关机方法
	req := &pb.VMRequest{Domain: inputParams.Domain}
	resp, err := instance.Running(context.Background(), req)
	if err != nil { //没启动或者找不到
		out := &dto.HostStatusOutput{Domain: inputParams.Domain, HostIP: inputParams.HostIP, Errmsg: fmt.Sprint(err)}
		middleware.ResponseSuccess(c, out)
		return
	}
	// ------------------------------完
	out := &dto.HostStatusOutput{Domain: inputParams.Domain, HostIP: inputParams.HostIP, Status: resp.Status, Errmsg: resp.Errmsg, IsRunning: resp.Running}
	middleware.ResponseSuccess(c, out)
}

// RecoverDoamin godoc
// @Summary 还原域
// @Description 还原域
// @Tags 主机管理
// @ID /host/control/recovery
// @Accept json
// @Produce json
// @Param domain query string true "域名"
// @Param host_ip query string true "宿主机ip"
// @Success 200 {object} middleware.Response{data=dto.HostStatusOutput} "success"
// @Router /host/control/recover [put]
func (s *HostController) RecoverDoamin(c *gin.Context) {
	inputParams := &dto.HostInfoInput{}
	if err := inputParams.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	// + GRPC 调用 ----------------------
	// 连接
	srvAddr := fmt.Sprintf("%s:%v", inputParams.HostIP, 50061)
	conn, err := grpc.Dial(srvAddr, grpc.WithInsecure())
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	defer conn.Close()

	// 初始化客户端
	instance := pb.NewVMClient(conn)
	// c := pb.NewVMClient(conn)
	// 调用关机方法
	req := &pb.VMRequest{Domain: inputParams.Domain}
	resp, err := instance.Recovery(context.Background(), req)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	// ------------------------------完
	out := &dto.HostStatusOutput{Domain: inputParams.Domain, HostIP: inputParams.HostIP, Status: resp.Status, Errmsg: resp.Errmsg, IsRunning: resp.Running}
	middleware.ResponseSuccess(c, out)
}
