package controller

import (
	"fmt"
	"strings"
	"time"

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
	group.GET("/service", service.ServiceDetail)      //获取单个服务详情
	group.GET("/statistic", service.ServiceStatistic) //数据统计
	group.DELETE("/service", service.ServiceDelete)
	group.POST("/service/http", service.ServiceAddHTTP)
	group.PUT("/service/http", service.ServiceUpdateHTTP)

	group.POST("/tcp", service.ServiceAddTCP)
	group.PUT("/tcp", service.ServiceUpdateTCP)
	group.POST("/grpc", service.ServiceAddGrpc)
	group.PUT("/grpc", service.ServiceUpdateGrpc)
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

// ServiceDetail godoc
// @Summary 服务详情
// @Description 服务详情
// @Tags 服务管理
// @ID /service/service/get
// @Accept json
// @Produce json
// @Param id query string true "服务ID"
// @Success 200 {object} middleware.Response{data=dao.ServiceDetail} "success"
// @Router /service/service [get]
func (s *ServiceController) ServiceDetail(c *gin.Context) {
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
	sdetail, err := serviceInfo.GetServiceDetail(c, tx, serviceInfo)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}

	out := sdetail
	middleware.ResponseSuccess(c, out)
}

// ServiceDelete godoc
// @Summary 服务删除
// @Description 服务删除
// @Tags 服务管理
// @ID /service/service/del
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
	//2. 判断ip列表和权重列表是否数量一致
	if len(strings.Split(inputParams.IPList, "\n")) != len(strings.Split(inputParams.WeightList, "\n")) {
		middleware.ResponseError(c, 2001, errors.New("ip列表和权重列表数量不一致"))
		return
	}
	//3. 从DB读取服务信息，判断服务名是否已存在
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	tx = tx.Begin() //开始事务
	serviceInfo := &dao.ServiceInfo{ServiceName: inputParams.ServiceName}
	if _, err = serviceInfo.Find(c, tx, serviceInfo); err == nil {
		tx.Rollback() //事务回滚
		middleware.ResponseError(c, 2003, errors.New("服务已存在"))
		return
	}
	//4. 从DB读取服务http url规则，判断是否已经在使用
	httpURL := &dao.HTTPRule{RuleType: inputParams.RuleType, Rule: inputParams.Rule}
	if _, err = httpURL.Find(c, tx, httpURL); err == nil {
		tx.Rollback() //事务回滚
		middleware.ResponseError(c, 2004, errors.New("服务接入前缀或域名已存在"))
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

// ServiceUpdateHTTP godoc
// @Summary 修改HTTP服务
// @Description 修改HTTP服务
// @Tags 服务管理
// @ID /service/service/http/put
// @Accept json
// @Produce json
// @Param body body dto.ServiceUpdateHTTPInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/service/http [put]
func (s *ServiceController) ServiceUpdateHTTP(c *gin.Context) {
	//1. 请求参数初步校验(必填)
	inputParams := &dto.ServiceUpdateHTTPInput{}
	if err := inputParams.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	if len(strings.Split(inputParams.IPList, "\n")) != len(strings.Split(inputParams.WeightList, "\n")) {
		middleware.ResponseError(c, 2001, errors.New("ip列表和权重列表数量不一致"))
		return
	}
	//2. 创建据库连接池
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	//3. 从DB从找服务基本信息
	serviceInfo := &dao.ServiceInfo{ServiceName: inputParams.ServiceName}
	serviceInfo, err = serviceInfo.Find(c, tx, serviceInfo)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	//4. 从DB中关联查询得到所有服务信息
	serviceDetail, err := serviceInfo.GetServiceDetail(c, tx, serviceInfo)
	if err != nil {
		middleware.ResponseError(c, 2004, err)
		return
	}
	tx = tx.Begin() //开始事务
	//5. 更新基本信息表
	//？？？根据名字查，不会重名吗
	srvInfo := serviceDetail.Info
	srvInfo.ServiceDesc = inputParams.ServiceDesc
	if srvInfo.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2005, err)
		return
	}
	//4. 更新httprule表(对指针的部分数据进行更新)
	httpRule := serviceDetail.HTTPRule
	httpRule.NeedHTTPS = inputParams.NeedHTTPS
	httpRule.NeedStripURI = inputParams.NeedStripURI
	httpRule.NeedWebsocket = inputParams.NeedWebsocket
	httpRule.URLRewrite = inputParams.HeaderTransfor
	httpRule.HeaderTransfor = inputParams.HeaderTransfor
	if httpRule.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2006, err)
		return
	}
	//5. 更新accessControl表(权限控制表)
	accessControl := serviceDetail.AccessControl
	accessControl.OpenAuth = inputParams.OpenAuth
	accessControl.BlackList = inputParams.BlackList
	accessControl.WhiteList = inputParams.WhiteList
	accessControl.ClientipFlowLimit = inputParams.ClientIPFlowLimit
	accessControl.ServiceFlowLimit = inputParams.ServiceFlowLimit
	if accessControl.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2007, err)
		return
	}
	//6. 更新loadbalance表(负载表)
	loadbalance := serviceDetail.LoadBalance
	loadbalance.RoundType = inputParams.RoundType
	loadbalance.IPList = inputParams.IPList
	loadbalance.WeightList = inputParams.WeightList
	loadbalance.UpstreamConnectTimeout = inputParams.UpstreamConnectTimeout
	loadbalance.UpstreamHeaderTimeout = inputParams.UpstreamHeaderTimeout
	loadbalance.UpstreamIdleTimeout = inputParams.UpstreamIdleTimeout
	loadbalance.UpstreamMaxIdle = inputParams.UpstreamMaxIdle
	if loadbalance.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2008, err)
		return
	}
	tx.Commit()
	//7. 返回信息
	out := fmt.Sprintf("create service [%v] success.", inputParams.ServiceName)
	middleware.ResponseSuccess(c, out)
}

// ServiceStatistic godoc
// @Summary 服务统计
// @Description 服务统计
// @Tags 服务管理
// @ID /service/statistic/get
// @Accept json
// @Produce json
// @Param id query string true "服务ID"
// @Success 200 {object} middleware.Response{data=dto.ServiceStatOutput} "success"
// @Router /service/statistic [get]
func (s *ServiceController) ServiceStatistic(c *gin.Context) {
	params := &dto.ServiceDeleteInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	// tx, err := lib.GetGormPool("default")
	// if err != nil {
	// 	middleware.ResponseError(c, 2001, err)
	// 	return
	// }
	// //从db中读取基本信息
	// serviceInfo := &dao.ServiceInfo{ID: params.ID}
	// serviceInfo, err = serviceInfo.Find(c, tx, serviceInfo)
	// if err != nil {
	// 	middleware.ResponseError(c, 2002, err)
	// 	return
	// }
	// sdetail, err := serviceInfo.GetServiceDetail(c, tx, serviceInfo)
	// if err != nil {
	// 	middleware.ResponseError(c, 2003, err)
	// 	return
	// }

	todayList := []int64{}
	for i := 0; i < time.Now().Hour(); i++ {
		todayList = append(todayList, 0)
	}
	yesterdayList := make([]int64, 24)

	out := dto.ServiceStatOutput{Today: todayList, Yesterday: yesterdayList}
	middleware.ResponseSuccess(c, out)
}

// ServiceAddTCP godoc
// @Summary tcp服务添加
// @Description tcp服务添加
// @Tags 服务管理
// @ID /tcp/add
// @Accept  json
// @Produce  json
// @Param body body dto.ServiceAddTCPInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /tcp [post]
func (s *ServiceController) ServiceAddTCP(c *gin.Context) {
	params := &dto.ServiceAddTCPInput{}
	if err := params.GetValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	//验证 service_name 是否被占用
	infoSearch := &dao.ServiceInfo{
		ServiceName: params.ServiceName,
		IsDelete:    0,
	}
	if _, err := infoSearch.Find(c, lib.GORMDefaultPool, infoSearch); err == nil {
		middleware.ResponseError(c, 2002, errors.New("服务名被占用，请重新输入"))
		return
	}

	//验证端口是否被占用?
	tcpRuleSearch := &dao.TCPRule{
		Port: params.Port,
	}
	if _, err := tcpRuleSearch.Find(c, lib.GORMDefaultPool, tcpRuleSearch); err == nil {
		middleware.ResponseError(c, 2003, errors.New("服务端口被占用，请重新输入"))
		return
	}
	grpcRuleSearch := &dao.GrpcRule{
		Port: params.Port,
	}
	if _, err := grpcRuleSearch.Find(c, lib.GORMDefaultPool, grpcRuleSearch); err == nil {
		middleware.ResponseError(c, 2004, errors.New("服务端口被占用，请重新输入"))
		return
	}

	//ip与权重数量一致
	if len(strings.Split(params.IPList, ",")) != len(strings.Split(params.WeightList, ",")) {
		middleware.ResponseError(c, 2005, errors.New("ip列表与权重设置不匹配"))
		return
	}

	tx := lib.GORMDefaultPool.Begin()
	info := &dao.ServiceInfo{
		LoadType:    public.LoadTypeTCP,
		ServiceName: params.ServiceName,
		ServiceDesc: params.ServiceDesc,
	}
	if err := info.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2006, err)
		return
	}
	loadBalance := &dao.LoadBalance{
		ServiceID:  info.ID,
		RoundType:  params.RoundType,
		IPList:     params.IPList,
		WeightList: params.WeightList,
		ForbidList: params.ForbidList,
	}
	if err := loadBalance.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2007, err)
		return
	}

	httpRule := &dao.TCPRule{
		ServiceID: info.ID,
		Port:      params.Port,
	}
	if err := httpRule.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2008, err)
		return
	}

	accessControl := &dao.AccessControl{
		ServiceID:         info.ID,
		OpenAuth:          params.OpenAuth,
		BlackList:         params.BlackList,
		WhiteList:         params.WhiteList,
		WhiteHostName:     params.WhiteHostName,
		ClientipFlowLimit: params.ClientIPFlowLimit,
		ServiceFlowLimit:  params.ServiceFlowLimit,
	}
	if err := accessControl.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2009, err)
		return
	}
	tx.Commit()
	middleware.ResponseSuccess(c, "")
	return
}

// ServiceUpdateTCP godoc
// @Summary tcp服务更新
// @Description tcp服务更新
// @Tags 服务管理
// @ID /tcp/update
// @Accept  json
// @Produce  json
// @Param body body dto.ServiceUpdateTCPInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /tcp [put]
func (s *ServiceController) ServiceUpdateTCP(c *gin.Context) {
	params := &dto.ServiceUpdateTCPInput{}
	if err := params.GetValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	//ip与权重数量一致
	if len(strings.Split(params.IPList, ",")) != len(strings.Split(params.WeightList, ",")) {
		middleware.ResponseError(c, 2002, errors.New("ip列表与权重设置不匹配"))
		return
	}

	tx := lib.GORMDefaultPool.Begin()

	service := &dao.ServiceInfo{
		ID: params.ID,
	}
	detail, err := service.GetServiceDetail(c, lib.GORMDefaultPool, service)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	info := detail.Info
	info.ServiceDesc = params.ServiceDesc
	if err := info.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2003, err)
		return
	}

	loadBalance := &dao.LoadBalance{}
	if detail.LoadBalance != nil {
		loadBalance = detail.LoadBalance
	}
	loadBalance.ServiceID = info.ID
	loadBalance.RoundType = params.RoundType
	loadBalance.IPList = params.IPList
	loadBalance.WeightList = params.WeightList
	loadBalance.ForbidList = params.ForbidList
	if err := loadBalance.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2004, err)
		return
	}

	tcpRule := &dao.TCPRule{}
	if detail.TCPRule != nil {
		tcpRule = detail.TCPRule
	}
	tcpRule.ServiceID = info.ID
	tcpRule.Port = params.Port
	if err := tcpRule.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2005, err)
		return
	}

	accessControl := &dao.AccessControl{}
	if detail.AccessControl != nil {
		accessControl = detail.AccessControl
	}
	accessControl.ServiceID = info.ID
	accessControl.OpenAuth = params.OpenAuth
	accessControl.BlackList = params.BlackList
	accessControl.WhiteList = params.WhiteList
	accessControl.WhiteHostName = params.WhiteHostName
	accessControl.ClientipFlowLimit = params.ClientIPFlowLimit
	accessControl.ServiceFlowLimit = params.ServiceFlowLimit
	if err := accessControl.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2006, err)
		return
	}
	tx.Commit()
	middleware.ResponseSuccess(c, "")
	return
}

// ServiceAddGrpc godoc
// @Summary grpc服务添加
// @Description grpc服务添加
// @Tags 服务管理
// @ID /grpc/add
// @Accept  json
// @Produce  json
// @Param body body dto.ServiceAddGrpcInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /grpc [post]
func (s *ServiceController) ServiceAddGrpc(c *gin.Context) {
	params := &dto.ServiceAddGrpcInput{}
	if err := params.GetValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	//验证 service_name 是否被占用
	infoSearch := &dao.ServiceInfo{
		ServiceName: params.ServiceName,
		IsDelete:    0,
	}
	if _, err := infoSearch.Find(c, lib.GORMDefaultPool, infoSearch); err == nil {
		middleware.ResponseError(c, 2002, errors.New("服务名被占用，请重新输入"))
		return
	}

	//验证端口是否被占用?
	tcpRuleSearch := &dao.TCPRule{
		Port: params.Port,
	}
	if _, err := tcpRuleSearch.Find(c, lib.GORMDefaultPool, tcpRuleSearch); err == nil {
		middleware.ResponseError(c, 2003, errors.New("服务端口被占用，请重新输入"))
		return
	}
	grpcRuleSearch := &dao.GrpcRule{
		Port: params.Port,
	}
	if _, err := grpcRuleSearch.Find(c, lib.GORMDefaultPool, grpcRuleSearch); err == nil {
		middleware.ResponseError(c, 2004, errors.New("服务端口被占用，请重新输入"))
		return
	}

	//ip与权重数量一致
	if len(strings.Split(params.IPList, ",")) != len(strings.Split(params.WeightList, ",")) {
		middleware.ResponseError(c, 2005, errors.New("ip列表与权重设置不匹配"))
		return
	}

	tx := lib.GORMDefaultPool.Begin()
	info := &dao.ServiceInfo{
		LoadType:    public.LoadTypeGRPC,
		ServiceName: params.ServiceName,
		ServiceDesc: params.ServiceDesc,
	}
	if err := info.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2006, err)
		return
	}

	loadBalance := &dao.LoadBalance{
		ServiceID:  info.ID,
		RoundType:  params.RoundType,
		IPList:     params.IPList,
		WeightList: params.WeightList,
		ForbidList: params.ForbidList,
	}
	if err := loadBalance.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2007, err)
		return
	}

	grpcRule := &dao.GrpcRule{
		ServiceID:      info.ID,
		Port:           params.Port,
		HeaderTransfor: params.HeaderTransfor,
	}
	if err := grpcRule.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2008, err)
		return
	}

	accessControl := &dao.AccessControl{
		ServiceID:         info.ID,
		OpenAuth:          params.OpenAuth,
		BlackList:         params.BlackList,
		WhiteList:         params.WhiteList,
		WhiteHostName:     params.WhiteHostName,
		ClientipFlowLimit: params.ClientIPFlowLimit,
		ServiceFlowLimit:  params.ServiceFlowLimit,
	}
	if err := accessControl.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2009, err)
		return
	}
	tx.Commit()
	middleware.ResponseSuccess(c, "")
	return
}

// ServiceUpdateGrpc godoc
// @Summary grpc服务更新
// @Description grpc服务更新
// @Tags 服务管理
// @ID /grpc/update
// @Accept  json
// @Produce  json
// @Param body body dto.ServiceUpdateGrpcInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /grpc [put]
func (s *ServiceController) ServiceUpdateGrpc(c *gin.Context) {
	params := &dto.ServiceUpdateGrpcInput{}
	if err := params.GetValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	//ip与权重数量一致
	if len(strings.Split(params.IPList, ",")) != len(strings.Split(params.WeightList, ",")) {
		middleware.ResponseError(c, 2002, errors.New("ip列表与权重设置不匹配"))
		return
	}

	tx := lib.GORMDefaultPool.Begin()

	service := &dao.ServiceInfo{
		ID: params.ID,
	}
	detail, err := service.GetServiceDetail(c, lib.GORMDefaultPool, service)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}

	info := detail.Info
	info.ServiceDesc = params.ServiceDesc
	if err := info.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2004, err)
		return
	}

	loadBalance := &dao.LoadBalance{}
	if detail.LoadBalance != nil {
		loadBalance = detail.LoadBalance
	}
	loadBalance.ServiceID = info.ID
	loadBalance.RoundType = params.RoundType
	loadBalance.IPList = params.IPList
	loadBalance.WeightList = params.WeightList
	loadBalance.ForbidList = params.ForbidList
	if err := loadBalance.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2005, err)
		return
	}

	grpcRule := &dao.GrpcRule{}
	if detail.GRPCRule != nil {
		grpcRule = detail.GRPCRule
	}
	grpcRule.ServiceID = info.ID
	//grpcRule.Port = params.Port
	grpcRule.HeaderTransfor = params.HeaderTransfor
	if err := grpcRule.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2006, err)
		return
	}

	accessControl := &dao.AccessControl{}
	if detail.AccessControl != nil {
		accessControl = detail.AccessControl
	}
	accessControl.ServiceID = info.ID
	accessControl.OpenAuth = params.OpenAuth
	accessControl.BlackList = params.BlackList
	accessControl.WhiteList = params.WhiteList
	accessControl.WhiteHostName = params.WhiteHostName
	accessControl.ClientipFlowLimit = params.ClientIPFlowLimit
	accessControl.ServiceFlowLimit = params.ServiceFlowLimit
	if err := accessControl.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2007, err)
		return
	}
	tx.Commit()
	middleware.ResponseSuccess(c, "")
	return
}
