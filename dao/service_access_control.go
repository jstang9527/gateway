package dao

import (
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"github.com/jstang9527/gateway/public"
)

// AccessControl ...
type AccessControl struct {
	ID                int64  `json:"id" gorm:"primary_key"`
	ServiceID         int64  `json:"service_id" gorm:"service_id" description:"服务id"`
	OpenAuth          int    `json:"open_auth" gorm:"open_auth" description:"是否开启权限 1=开启"`
	BlackList         string `json:"black_list" gorm:"black_list" description:"黑名单ip"`
	WhiteList         string `json:"white_list" gorm:"white_list" description:"白名单ip"`
	WhiteHostName     string `json:"white_host_name" gorm:"white_host_name" description:"白名单主机"`
	ClientipFlowLimit int    `json:"clientip_flow_limit" gorm:"clientip_flow_limit" description:"客户端ip限流"`
	ServiceFlowLimit  int    `json:"service_flow_limit" gorm:"service_flow_limit" description:"服务端限流"`
}

// TableName ...
func (a *AccessControl) TableName() string {
	return "gateway_service_access_control"
}

// Find ...
func (a *AccessControl) Find(c *gin.Context, tx *gorm.DB, search *AccessControl) (*AccessControl, error) {
	model := &AccessControl{}
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(model).Error
	return model, err
}

// Save ...
func (a *AccessControl) Save(c *gin.Context, tx *gorm.DB) error {
	if err := tx.SetCtx(public.GetGinTraceContext(c)).Save(a).Error; err != nil {
		return err
	}
	return nil
}
