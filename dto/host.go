package dto

import (
	"github.com/gin-gonic/gin"
	"github.com/jstang9527/gateway/public"
)

// HostListInput 主机分页查询
type HostListInput struct {
	Info     string `json:"info" form:"info" comment:"关键词" example:"" validate:""`                      //关键词
	PageNo   int    `json:"page_no" form:"page_no" comment:"页数" example:"1" validate:"required"`        //页数
	PageSize int    `json:"page_size" form:"page_size" comment:"每页条数" example:"20" validate:"required"` // 每页条数
}

// HostListOutput ...
type HostListOutput struct {
	Total int64                `json:"total" form:"total" comment:"总数" example:"" validate:""` //总数
	List  []HostListItemOutput `json:"list" form:"list" comment:"列表" example:"" validate:""`   //列表
}

// BindValidParam 校验删除参数,绑定结构体,校验参数
func (s *HostListInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, s)
}

// HostListItemOutput 域item信息
type HostListItemOutput struct {
	ID           int64  `json:"id" form:"id"`                       //id
	Domain       string `json:"domain" form:"domain"`               //主机名
	DomainIP     string `json:"domain_ip" form:"domain_ip"`         //主机ip
	DomainOS     string `json:"domain_os" form:"domain_os"`         //操作系统
	DomainType   int    `json:"domain_type" form:"domain_type"`     //主机类型:0真机,1云机,2虚拟机
	DomainDesc   string `json:"domain_desc" form:"domain_desc"`     //主机描述
	HostIP       string `json:"host_ip" form:"host_ip"`             //宿主机ip
	HostDesc     string `json:"host_desc" form:"host_desc"`         //宿主机描述
	DomainStatus bool   `json:"domain_status" form:"domain_status"` //占位字段
}
