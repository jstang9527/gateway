package dao

import (
	"time"

	"github.com/jstang9527/gateway/dto"

	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"github.com/jstang9527/gateway/public"
)

// ServiceInfo 服务信息
type ServiceInfo struct {
	ID          int64     `json:"id" gorm:"primary_key"`
	LoadType    int       `json:"load_type" gorm:"column:load_type" description:"负载类型 0=http 1=tcp 2=grpc"`
	ServiceName string    `json:"service_name" gorm:"column:service_name" description:"服务名"`
	ServiceDesc string    `json:"service_desc" gorm:"column:service_desc" description:"服务描述"`
	UpdatedAt   time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	CreatedAt   time.Time `json:"create_at" gorm:"column:create_at" description:"创建时间"`
	IsDelete    int8      `json:"is_delete" gorm:"column:is_delete" description:"是否删除;0:否;1:是"`
}

// TableName 表名
func (s *ServiceInfo) TableName() string {
	return "gateway_service_info"
}

// ServiceDetail 服务详情 service_info + service_rule
// 通过serviceInfo的ID查service_rule
// 得到的结果封装到完全体
func (s *ServiceInfo) ServiceDetail(c *gin.Context, tx *gorm.DB, search *ServiceInfo) (*ServiceDetail, error) {
	httpRule := &HTTPRule{ServiceID: search.ID}
	httpRule, err := httpRule.Find(c, tx, httpRule)
	if err != nil && err != gorm.ErrRecordNotFound { //找不到也是正确的
		return nil, err
	}

	tcpRule := &TCPRule{ServiceID: search.ID}
	tcpRule, err = tcpRule.Find(c, tx, tcpRule)
	if err != nil && err != gorm.ErrRecordNotFound { //找不到也是正确的
		return nil, err
	}

	grpcRule := &GrpcRule{ServiceID: search.ID}
	grpcRule, err = grpcRule.Find(c, tx, grpcRule)
	if err != nil && err != gorm.ErrRecordNotFound { //找不到也是正确的
		return nil, err
	}

	accessControl := &AccessControl{ServiceID: search.ID}
	accessControl, err = accessControl.Find(c, tx, accessControl)
	if err != nil && err != gorm.ErrRecordNotFound { //找不到也是正确的
		return nil, err
	}

	loadBalance := &LoadBalance{ServiceID: search.ID}
	loadBalance, err = loadBalance.Find(c, tx, loadBalance)
	if err != nil && err != gorm.ErrRecordNotFound { //找不到也是正确的
		return nil, err
	}

	detail := &ServiceDetail{
		Info:          search,
		HTTPRule:      httpRule,
		TCPRule:       tcpRule,
		GRPCRule:      grpcRule,
		LoadBalance:   loadBalance,
		AccessControl: accessControl,
	}
	return detail, nil
}

// Find 查服务gateway_info
func (s *ServiceInfo) Find(c *gin.Context, tx *gorm.DB, search *ServiceInfo) (*ServiceInfo, error) {
	out := &ServiceInfo{}
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Save 数据库表修改
func (s *ServiceInfo) Save(c *gin.Context, tx *gorm.DB) error {
	return tx.SetCtx(public.GetGinTraceContext(c)).Save(s).Error
}

// PageList 分页查询=>得服务数组
func (s *ServiceInfo) PageList(c *gin.Context, tx *gorm.DB, search *dto.ServiceListInput) ([]ServiceInfo, int64, error) {
	var total int64 = 0
	list := []ServiceInfo{}
	offset := (search.PageNo - 1) * search.PageSize
	query := tx.SetCtx(public.GetGinTraceContext(c))
	query = query.Table(s.TableName()).Where("is_delete=0")
	if search.Info != "" {
		query = query.Where("service_name like ? or service_desc like ?", "%"+search.Info+"%", "%"+search.Info+"%")
	}
	if err := query.Limit(search.PageSize).Offset(offset).Order("id desc").Find(&list).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	query.Limit(search.PageSize).Offset(offset).Count(&total)
	return list, total, nil
}
