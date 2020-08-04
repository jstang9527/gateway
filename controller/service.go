package controller

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

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
	// group.POST("/service_add_http", service.ServiceAddHTTP)
	group.POST("/service/http", service.ServiceAddHTTP)
	group.PUT("/service/http", service.ServiceAddHTTP)
}

// ServiceList godoc
// @Summary 服务列表
// @Description 服务列表
// @Tags 服务管理
// @ID /service/service_list
// @Accept json
// @Produce json
// @Param info query string false "关键词"
// @Param page_size query int true "每页个数"
// @Param page_no query int true "当前页数"
// @Success 200 {object} middleware.Response{data=dto.ServiceListOutput} "success"
// @Router /service/service_list [get]
func (s *ServiceController) ServiceList(c *gin.Context) {
	inputParams := &dto.ServiceListInput{}
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
	serviceInfo := &dao.ServiceInfo{}
	serviceList, total, err := serviceInfo.PageList(c, tx, inputParams)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	//格式化输出信息
	outList := []dto.ServiceListItemOutput{} //这个结构体是面向前端接口的
	for _, item := range serviceList {
		servicedetail, err := item.GetServiceDetail(c, tx, &item)
		if err != nil {
			middleware.ResponseError(c, 2003, err)
			return
		}
		// 1.http后缀接入clusterIP+clusterPort+path
		// 2.http域名接入domain
		// 3.tcp、grpc接入clusterIP+servicePort
		serviceAddr := "Unknow"
		clusterIP := lib.GetStringConf("base.cluster.cluster_ip")
		clusterPort := lib.GetStringConf("base.cluster.cluster_port")
		clusterSslPort := lib.GetStringConf("base.cluster.cluster_ssl_port")
		if servicedetail.Info.LoadType == public.LoadTypeHTTP && servicedetail.HTTPRule.RuleType == public.HTTPRuleTypePrefixURL && servicedetail.HTTPRule.NeedHTTPS == 1 {
			serviceAddr = fmt.Sprintf("%s:%s%s", clusterIP, clusterSslPort, servicedetail.HTTPRule.Rule)
		}
		if servicedetail.Info.LoadType == public.LoadTypeHTTP && servicedetail.HTTPRule.RuleType == public.HTTPRuleTypePrefixURL && servicedetail.HTTPRule.NeedHTTPS == 0 {
			serviceAddr = fmt.Sprintf("%s:%s%s", clusterIP, clusterPort, servicedetail.HTTPRule.Rule)
		}
		if servicedetail.Info.LoadType == public.LoadTypeHTTP && servicedetail.HTTPRule.RuleType == public.HTTPRuleTypeDomain {
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
			ID:          item.ID,
			ServiceName: item.ServiceName,
			ServiceDesc: item.ServiceDesc,
			ServiceAddr: serviceAddr,
			LoadType:    item.LoadType,
			QPS:         0,
			QPD:         0,
			TotalNode:   len(ipList),
		}
		outList = append(outList, outItem)
	}
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

// ServiceAddHTTP godoc
// @Summary 添加HTTP服务
// @Description 添加HTTP服务
// @Tags 服务管理
// @ID /service/service/http/post
// @Accept json
// @Produce json
// @Param body body dto.ServiceAddHTTPInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/service/http [post]
func (s *ServiceController) ServiceAddHTTP(c *gin.Context) {
	//1. 请求参数初步校验(必填)
	inputParams := &dto.ServiceAddHTTPInput{}
	if err := inputParams.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	//2. 从DB读取服务信息，判断服务名是否已存在
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	tx = tx.Begin() //开始事务
	serviceInfo := &dao.ServiceInfo{ServiceName: inputParams.ServiceName}
	if _, err = serviceInfo.Find(c, tx, serviceInfo); err == nil {
		tx.Rollback() //事务回滚
		middleware.ResponseError(c, 2002, errors.New("服务已存在"))
		return
	}
	//3. 从DB读取服务http url规则，判断是否已经在使用
	httpURL := &dao.HTTPRule{RuleType: inputParams.RuleType, Rule: inputParams.Rule}
	if _, err = httpURL.Find(c, tx, httpURL); err == nil {
		tx.Rollback() //事务回滚
		middleware.ResponseError(c, 2003, errors.New("服务接入前缀或域名已存在"))
		return
	}
	//4. 判断ip列表和权重列表是否数量一致
	if len(strings.Split(inputParams.IPList, "\n")) != len(strings.Split(inputParams.WeightList, "\n")) {
		tx.Rollback() //事务回滚
		middleware.ResponseError(c, 2004, errors.New("ip列表和权重列表数量不一致"))
		return
	}
	//5. 入库
	//5.1 基本服务信息表
	serviceModel := &dao.ServiceInfo{
		ServiceName: inputParams.ServiceName,
		ServiceDesc: inputParams.ServiceDesc}
	if err := serviceModel.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2005, err)
		return
	}
	//5.2 http服务规则表
	// 拿主键，上面save完后serviceModel就会有id值
	fmt.Println("serviceModel_id ==>", serviceModel.ID)
	httpRuleModel := &dao.HTTPRule{
		ServiceID:      serviceModel.ID,
		RuleType:       inputParams.RuleType,
		Rule:           inputParams.Rule,
		NeedHTTPS:      inputParams.NeedHTTPS,
		NeedStripURI:   inputParams.NeedStripURI,
		NeedWebsocket:  inputParams.NeedWebsocket,
		URLRewrite:     inputParams.URLRewrite,
		HeaderTransfor: inputParams.HeaderTransfor,
	}
	if err := httpRuleModel.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2006, err)
		return
	}
	//5.3 权限表
	accessControlModel := &dao.AccessControl{
		ServiceID:         serviceModel.ID,
		OpenAuth:          inputParams.OpenAuth,
		BlackList:         inputParams.BlackList,
		WhiteList:         inputParams.WhiteList,
		ClientipFlowLimit: inputParams.ClientIPFlowLimit,
		ServiceFlowLimit:  inputParams.ServiceFlowLimit,
	}
	if err := accessControlModel.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2007, err)
		return
	}
	//5.4 负载均衡表
	loadbalanceModel := &dao.LoadBalance{
		ServiceID:              serviceModel.ID,
		RoundType:              inputParams.RoundType,
		IPList:                 inputParams.IPList,
		WeightList:             inputParams.WeightList,
		UpstreamConnectTimeout: inputParams.UpstreamConnectTimeout,
		UpstreamHeaderTimeout:  inputParams.UpstreamHeaderTimeout,
		UpstreamIdleTimeout:    inputParams.UpstreamIdleTimeout,
		UpstreamMaxIdle:        inputParams.UpstreamMaxIdle,
	}
	if err := loadbalanceModel.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2008, err)
		return
	}
	tx.Commit() //提交事务
	//6. 返回信息
	out := fmt.Sprintf("create service [%v] success.", inputParams.ServiceName)
	middleware.ResponseSuccess(c, out)
}

// ServicePutHTTP godoc
// @Summary 修改HTTP服务
// @Description 修改HTTP服务
// @Tags 服务管理
// @ID /service/service/http/put
// @Accept json
// @Produce json
// @Param body body dto.ServiceAddHTTPInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/service/http [put]
func (s *ServiceController) ServicePutHTTP(c *gin.Context) {
	//1. 请求参数初步校验(必填)
	inputParams := &dto.ServiceAddHTTPInput{}
	if err := inputParams.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	//6. 返回信息
	out := fmt.Sprintf("create service [%v] success.", inputParams.ServiceName)
	middleware.ResponseSuccess(c, out)
}
