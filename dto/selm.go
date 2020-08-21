package dto

import (
	"github.com/gin-gonic/gin"
	"github.com/jstang9527/gateway/public"
)

// WebhookInput 主机分页查询
type WebhookInput struct {
	ProjectName   string `json:"project_name" form:"project_name" comment:"项目名" example:"applet" validate:"required"`              //项目名
	ProjectAddr   string `json:"project_addr" form:"project_addr" comment:"项目地址" example:"https://gitlab.com" validate:""`         //项目地址
	WebAddr       string `json:"web_addr" form:"web_addr" comment:"Web地址" example:"http://172.31.50.39:65000" validate:"required"` // Web地址
	SearchTimeout int    `json:"search_timeout" form:"search_timeout" comment:"元素检索超时时间" example:"10" validate:""`                 //元素检索超时时间
	SyncNum       int    `json:"sync_num" form:"sync_num" comment:"并发数" example:"1" validate:""`                                   //并发数
	StreamID      int64  `json:"stream_id" form:"stream_id" comment:"流水线ID" example:"1" validate:"required"`                       //流水线ID
}

// BindValidParam 校验删除参数,绑定结构体,校验参数
func (s *WebhookInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, s)
}
