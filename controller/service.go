package controller

import (
	"fmt"

	"github.com/jstang9527/gateway/public"

	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"github.com/jstang9527/gateway/dao"
	"github.com/jstang9527/gateway/dto"
	"github.com/jstang9527/gateway/middleware"
)

// ServiceController ...
type ServiceController struct{}

// ServiceRegister ...
func ServiceRegister(group *gin.RouterGroup) {
	service := &ServiceController{}
	group.GET("/service_list", service.ServiceList)
	group.DELETE("/service", service.ServiceDelete)
}

// ServiceList godoc
// @Summary 服务列表
// @Description 服务列表
// @Tags 服务管理
// @ID /service/service_list
// @Accept json
// @Produce json
// @Param info query string false "关键词"
// @Param page_size query int true "每个个数"
// @Param page_no query int true "当前页数"
// @Success 200 {object} middleware.Response{data=dto.ServiceListOutput} "success"
// @Router /service/service_list [get]
func (s *ServiceController) ServiceList(c *gin.Context) {
	params := &dto.ServiceListInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	//从db中分页读取基本信息
	serviceInfo := &dao.ServiceInfo{}
	list, total, err := serviceInfo.PageList(c, tx, params)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	//格式化输出信息
	outList := []dto.ServiceListItemOutput{}
	for _, listitem := range list {
		servicedetail, err := listitem.ServiceDetail(c, tx, &listitem)
		if err != nil {
			middleware.ResponseError(c, 2003, err)
			return
		}
		// 1.http后缀接入clusterIP+clusterPort+path
		// 2.http域名接入domain
		// 3.tcp、grpc接入clusterIP+servicePort
		serviceAddr := "unknow"
		clusterIP := lib.GetStringConf("base.cluster.cluster_ip")
		clusterPort := lib.GetStringConf("base.cluster.cluster_port")
		clusterSslPort := lib.GetStringConf("base.cluster.cluster_ssl_port")
		if servicedetail.Info.LoadType == public.LoadTypeHTTP &&
			servicedetail.HTTPRule.RuleType == public.HTTPRuleTypePrefixURL &&
			servicedetail.HTTPRule.NeedHTTPS == 1 {
			serviceAddr = fmt.Sprintf("%s:%s%s", clusterIP, clusterSslPort, servicedetail.HTTPRule.Rule)
		}
		if servicedetail.Info.LoadType == public.LoadTypeHTTP &&
			servicedetail.HTTPRule.RuleType == public.HTTPRuleTypePrefixURL &&
			servicedetail.HTTPRule.NeedHTTPS == 0 {
			serviceAddr = fmt.Sprintf("%s:%s%s", clusterIP, clusterPort, servicedetail.HTTPRule.Rule)
		}
		if servicedetail.Info.LoadType == public.LoadTypeHTTP &&
			servicedetail.HTTPRule.RuleType == public.HTTPRuleTypeDomain {
			serviceAddr = servicedetail.HTTPRule.Rule
		}
		if servicedetail.Info.LoadType == public.LoadTypeTCP {
			serviceAddr = fmt.Sprintf("%s:%v", clusterIP, servicedetail.TCPRule.Port)
		}
		if servicedetail.Info.LoadType == public.LoadTypeGRPC {
			serviceAddr = fmt.Sprintf("%s:%v", clusterIP, servicedetail.GRPCRule.Port)
		}
		ipList := servicedetail.LoadBalance.GetIPListByModel()
		outItem := dto.ServiceListItemOutput{
			ID:          listitem.ID,
			ServiceName: listitem.ServiceName,
			ServiceDesc: listitem.ServiceDesc,
			ServiceAddr: serviceAddr,
			QPS:         0,
			QPD:         0,
			TotalNode:   len(ipList),
		}
		outList = append(outList, outItem)
	}
	fmt.Println(total)
	out := &dto.ServiceListOutput{Total: total, List: outList}
	middleware.ResponseSuccess(c, out)
}

// ServiceDelete godoc
// @Summary 服务删除
// @Description 服务删除
// @Tags 服务管理
// @ID /service/service
// @Accept json
// @Produce json
// @Param id query string true "服务ID"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/service [delete]
func (s *ServiceController) ServiceDelete(c *gin.Context) {
	params := &dto.ServiceDeleteInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	//从db中读取基本信息
	serviceInfo := &dao.ServiceInfo{ID: params.ID}
	serviceInfo, err = serviceInfo.Find(c, tx, serviceInfo)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	serviceInfo.IsDelete = 1
	if err := serviceInfo.Save(c, tx); err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	out := fmt.Sprintf("delete success. id=%v", serviceInfo.ID)
	middleware.ResponseSuccess(c, out)
}
