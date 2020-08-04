package dto

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jstang9527/gateway/public"
)

// AdminLoginInput ...
type AdminLoginInput struct {
	UserName string `json:"username" form:"username" comment:"账户" example:"账户" validate:"required,valid_username"`
	Password string `json:"password" form:"password" comment:"密码" example:"密码" validate:"required"`
}

// AdminSessionInfo ...
type AdminSessionInfo struct {
	ID        int       `json:"id"`
	UserName  string    `json:"username"`
	LoginTime time.Time `json:"login_time"`
}

// BindValidParam 校验方法,绑定结构体,校验参数
func (a *AdminLoginInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, a)
}

// AdminLoginOutput ...
type AdminLoginOutput struct {
	Token string `json:"token" form:"token" comment:"token" example:"token" validate:""`
}
