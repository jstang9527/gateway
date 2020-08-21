package controller

import (
	"time"

	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"github.com/jstang9527/gateway/dao"
	"github.com/jstang9527/gateway/dto"
	"github.com/jstang9527/gateway/middleware"
	"github.com/jstang9527/gateway/public"
	"github.com/pkg/errors"
)

// AppController ...
type AppController struct{}

// AppRegister ...
func AppRegister(router *gin.RouterGroup) {
	admin := AppController{}
	router.GET("/list", admin.AppList)
	router.GET("/detail", admin.APPDetail)
	router.GET("/stat", admin.AppStatistics)
	router.DELETE("/app", admin.APPDelete)
	router.POST("/app", admin.AppAdd)
	router.PUT("/app", admin.AppUpdate)
}

// AppList godoc
// @Summary 项目列表
// @Description 项目列表
// @Tags 项目管理
// @ID /app/list
// @Accept  json
// @Produce  json
// @Param info query string false "关键词"
// @Param page_size query string true "每页多少条"
// @Param page_no query string true "页码"
// @Success 200 {object} middleware.Response{data=dto.APPListOutput} "success"
// @Router /app/list [get]
func (admin *AppController) AppList(c *gin.Context) {
	params := &dto.APPListInput{}
	if err := params.GetValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	info := &dao.App{}
	list, total, err := info.APPList(c, lib.GORMDefaultPool, params)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	outputList := []dto.APPListItemOutput{}
	for _, item := range list {
		var realQPS int64 = 0
		var realQpd int64 = 0
		outputList = append(outputList, dto.APPListItemOutput{
			ID:       item.ID,
			AppID:    item.AppID,
			Name:     item.Name,
			Secret:   item.Secret,
			WhiteIPS: item.WhiteIPS,
			Qpd:      item.Qpd,
			QPS:      item.QPS,
			RealQpd:  realQpd,
			RealQPS:  realQPS,
		})
	}
	output := dto.APPListOutput{
		List:  outputList,
		Total: total,
	}
	middleware.ResponseSuccess(c, output)
	return
}

// APPDetail godoc
// @Summary 项目详情
// @Description 项目详情
// @Tags 项目管理
// @ID /app/detail
// @Accept  json
// @Produce  json
// @Param id query string true "项目ID"
// @Success 200 {object} middleware.Response{data=dao.App} "success"
// @Router /app/detail [get]
func (admin *AppController) APPDetail(c *gin.Context) {
	params := &dto.APPDetailInput{}
	if err := params.GetValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	search := &dao.App{
		ID: params.ID,
	}
	detail, err := search.Find(c, lib.GORMDefaultPool, search)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	middleware.ResponseSuccess(c, detail)
	return
}

// APPDelete godoc
// @Summary 项目删除
// @Description 项目删除
// @Tags 项目管理
// @ID /app/delete
// @Accept  json
// @Produce  json
// @Param id query string true "项目ID"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /app/app [delete]
func (admin *AppController) APPDelete(c *gin.Context) {
	params := &dto.APPDetailInput{}
	if err := params.GetValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	search := &dao.App{
		ID: params.ID,
	}
	info, err := search.Find(c, lib.GORMDefaultPool, search)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	info.IsDelete = 1
	if err := info.Save(c, lib.GORMDefaultPool); err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	middleware.ResponseSuccess(c, "")
	return
}

// AppAdd godoc
// @Summary 项目添加
// @Description 项目添加
// @Tags 项目管理
// @ID /app/add
// @Accept  json
// @Produce  json
// @Param body body dto.APPAddHTTPInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /app/app [post]
func (admin *AppController) AppAdd(c *gin.Context) {
	params := &dto.APPAddHTTPInput{}
	if err := params.GetValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	//验证app_id是否被占用
	search := &dao.App{
		AppID: params.AppID,
	}
	if _, err := search.Find(c, lib.GORMDefaultPool, search); err == nil {
		middleware.ResponseError(c, 2002, errors.New("租户ID被占用，请重新输入"))
		return
	}
	if params.Secret == "" {
		params.Secret = public.MD5(params.AppID)
	}
	tx := lib.GORMDefaultPool
	info := &dao.App{
		AppID:    params.AppID,
		Name:     params.Name,
		Secret:   params.Secret,
		WhiteIPS: params.WhiteIPS,
		QPS:      params.QPS,
		Qpd:      params.Qpd,
	}
	if err := info.Save(c, tx); err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	middleware.ResponseSuccess(c, "")
	return
}

// AppUpdate godoc
// @Summary 项目更新
// @Description 项目更新
// @Tags 项目管理
// @ID /app/update
// @Accept  json
// @Produce  json
// @Param body body dto.APPUpdateHTTPInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /app/app [put]
func (admin *AppController) AppUpdate(c *gin.Context) {
	params := &dto.APPUpdateHTTPInput{}
	if err := params.GetValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	search := &dao.App{
		ID: params.ID,
	}
	info, err := search.Find(c, lib.GORMDefaultPool, search)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	if params.Secret == "" {
		params.Secret = public.MD5(params.AppID)
	}
	info.Name = params.Name
	info.Secret = params.Secret
	info.WhiteIPS = params.WhiteIPS
	info.QPS = params.QPS
	info.Qpd = params.Qpd
	if err := info.Save(c, lib.GORMDefaultPool); err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	middleware.ResponseSuccess(c, "")
	return
}

// AppStatistics godoc
// @Summary 项目统计
// @Description 项目统计
// @Tags 项目管理
// @ID /stat/get
// @Accept  json
// @Produce  json
// @Param id query string true "项目ID"
// @Success 200 {object} middleware.Response{data=dto.StatisticsOutput} "success"
// @Router /app/stat [get]
func (admin *AppController) AppStatistics(c *gin.Context) {
	params := &dto.APPDetailInput{}
	if err := params.GetValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	// search := &dao.App{
	// 	ID: params.ID,
	// }
	// detail, err := search.Find(c, lib.GORMDefaultPool, search)
	// if err != nil {
	// 	middleware.ResponseError(c, 2002, err)
	// 	return
	// }

	//今日流量全天小时级访问统计
	todayStat := []int64{}
	for i := 0; i <= time.Now().In(lib.TimeLocation).Hour(); i++ {
		todayStat = append(todayStat, 0)
	}

	//昨日流量全天小时级访问统计
	yesterdayStat := make([]int64, 24)

	stat := dto.StatisticsOutput{
		Today:     todayStat,
		Yesterday: yesterdayStat,
	}
	middleware.ResponseSuccess(c, stat)
	return
}
