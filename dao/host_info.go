package dao

import (
	"github.com/jstang9527/gateway/dto"

	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"github.com/jstang9527/gateway/public"
)

// HostInfo 服务信息
type HostInfo struct {
	ID         int64  `json:"id" gorm:"primary_key"`                                                 //id
	Domain     string `json:"domain" gorm:"column:domain" description:"主机名"`                         //主机名
	DomainIP   string `json:"domain_ip" gorm:"column:domain_ip" description:"主机ip"`                  //主机ip
	DomainOS   string `json:"domain_os" gorm:"column:domain_os" description:"操作系统"`                  //操作系统
	DomainType int    `json:"domain_type" gorm:"column:domain_type" description:"主机类型;0:真机1:云机2:虚机"` //主机类型:0真机,1云机,2虚拟机
	DomainDesc string `json:"domain_desc" gorm:"column:domain_desc" description:"主机描述"`              //主机描述
	HostIP     string `json:"host_ip" gorm:"column:host_ip" description:"宿主机ip"`                     //宿主机ip
	HostDesc   string `json:"host_desc" gorm:"column:host_desc" description:"宿主机描述"`                 //宿主机描述
}

// TableName 表名
func (s *HostInfo) TableName() string {
	return "host"
}

// Find 查服务host_info
func (s *HostInfo) Find(c *gin.Context, tx *gorm.DB, search *HostInfo) (*HostInfo, error) {
	out := &HostInfo{}
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Save 数据库表修改
func (s *HostInfo) Save(c *gin.Context, tx *gorm.DB) error {
	return tx.SetCtx(public.GetGinTraceContext(c)).Save(s).Error
}

// PageList 分页查询=>得服务数组
func (s *HostInfo) PageList(c *gin.Context, tx *gorm.DB, search *dto.HostListInput) ([]HostInfo, int64, error) {
	var total int64 = 0
	list := []HostInfo{}
	offset := (search.PageNo - 1) * search.PageSize //第一页的话就不用偏移了,直接limit查
	query := tx.SetCtx(public.GetGinTraceContext(c))
	// query = query.Table(s.TableName()).Where("is_delete=0")
	query = query.Table(s.TableName())
	if search.Info != "" { ///模糊查询
		query = query.Where("domain like ? or domain_desc like ?", "%"+search.Info+"%", "%"+search.Info+"%")
	}
	if err := query.Limit(search.PageSize).Offset(offset).Order("id desc").Find(&list).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	query.Limit(search.PageSize).Offset(offset).Count(&total)
	return list, total, nil
}
