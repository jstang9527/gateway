package dto

import (
	"github.com/gin-gonic/gin"
	"github.com/jstang9527/gateway/public"
)

// HostInfoInput 主机状态查询
type HostInfoInput struct {
	Domain string `json:"domain" form:"domain" comment:"域" example:"centos75" validate:"required"`            //域
	HostIP string `json:"host_ip" form:"host_ip" comment:"宿主机ip" example:"172.31.50.254" validate:"required"` //宿主机IP
}

// BindValidParam 校验删除参数,绑定结构体,校验参数
func (s *HostInfoInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, s)
}

// HostStatusOutput 域状态信息
type HostStatusOutput struct {
	Domain    string `json:"domain" form:"domain"`       //域
	HostIP    string `json:"host_ip" form:"host_ip"`     //宿主机名
	Status    bool   `json:"status" form:"status"`       //状态信息 bool
	Errmsg    string `json:"errmsg" form:"errmsg"`       //错误信息
	IsRunning bool   `json:"isRunning" form:"isRunning"` //状态信息 bool
}
