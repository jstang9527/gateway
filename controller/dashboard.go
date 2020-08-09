package controller

import (
	"errors"
	"time"

	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"github.com/jstang9527/gateway/dao"
	"github.com/jstang9527/gateway/dto"
	"github.com/jstang9527/gateway/middleware"
	"github.com/jstang9527/gateway/public"
)

// DashboardController ...
type DashboardController struct{}

// DashboardRegister ...
func DashboardRegister(group *gin.RouterGroup) {
	ctl := &DashboardController{}
	group.GET("/panel_data", ctl.PanelData)
	group.GET("/flow_stat", ctl.FlowStat)
	group.GET("/service_stat", ctl.ServiceStat)
}

// PanelData godoc
// @Summary 指标统计
// @Description 指标统计
// @Tags 首页大盘
// @ID /dashboard/panel_data/get
// @Accept json
// @Produce json
// @Success 200 {object} middleware.Response{data=dto.PanelDataOutput} "success"
// @Router /dashboard/panel_data [get]
func (s *DashboardController) PanelData(c *gin.Context) {
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	//1. 从DB查服务数量
	serviceInfo := &dao.ServiceInfo{}
	_, serviceNum, err := serviceInfo.PageList(c, tx, &dto.ServiceListInput{PageSize: 1, PageNo: 1})
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	//2. 从DB查租户数量
	app := &dao.App{}
	_, appNum, err := app.APPList(c, tx, &dto.APPListInput{PageNo: 1, PageSize: 1})
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	out := &dto.PanelDataOutput{
		ServiceNum:  serviceNum,
		AppNum:      appNum,
		TodayReqNum: 0,
		CurrentQPS:  0,
	}
	middleware.ResponseSuccess(c, out)
}

// FlowStat godoc
// @Summary 流量统计
// @Description 流量统计
// @Tags 首页大盘
// @ID /dashboard/flow_stat/get
// @Accept json
// @Produce json
// @Success 200 {object} middleware.Response{data=dto.ServiceStatOutput} "success"
// @Router /dashboard/flow_stat [get]
func (s *DashboardController) FlowStat(c *gin.Context) {
	todayList := []int64{}
	for i := 0; i < time.Now().Hour(); i++ {
		todayList = append(todayList, 0)
	}
	yesterdayList := make([]int64, 24)

	out := dto.ServiceStatOutput{Today: todayList, Yesterday: yesterdayList}
	middleware.ResponseSuccess(c, out)
}

// ServiceStat godoc
// @Summary 服务统计
// @Description 服务统计
// @Tags 首页大盘
// @ID /dashboard/service_stat/get
// @Accept json
// @Produce json
// @Success 200 {object} middleware.Response{data=dto.PanelSrvStatOutput} "success"
// @Router /dashboard/service_stat [get]
func (s *DashboardController) ServiceStat(c *gin.Context) {
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	// 1. 从DB查服务类型占比
	serviceInfo := &dao.ServiceInfo{}
	list, err := serviceInfo.GroupBySrvType(c, tx)
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	legend := []string{}
	for index, item := range list {
		name, ok := public.LoadTypeMap[item.LoadType]
		if !ok {
			middleware.ResponseError(c, 2002, errors.New("load_type not found"))
			return
		}
		list[index].Name = name
		legend = append(legend, name)
	}
	out := &dto.PanelSrvStatOutput{Legend: legend, Data: list}
	middleware.ResponseSuccess(c, out)
}
