package dto

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jstang9527/gateway/public"
)

// AdminInfoOutput 管理员信息输出
type AdminInfoOutput struct {
	ID           int       `json:"id"`
	UserName     string    `json:"user_name"`
	LoginTime    time.Time `json:"login_time"`
	Avatar       string    `json:"avatar"`
	Introduction string    `json:"introduction"`
	Roles        []string  `json:"roles"`
}

// ChangePwdInput ...
type ChangePwdInput struct {
	Password string `json:"password" form:"password" comment:"密码" example:"密码" validate:"required"`
}

// BindValidParam 校验方法,绑定结构体,校验参数
func (a *ChangePwdInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, a)
}
