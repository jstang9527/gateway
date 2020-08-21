package dao

import (
	"time"

	"github.com/jstang9527/gateway/dto"
	"github.com/jstang9527/gateway/public"

	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
)

// StreamInfo 组件由多个动作组成
type StreamInfo struct {
	ID         int64     `json:"id" gorm:"primary_key"`
	StreamName string    `json:"stream_name" gorm:"column:stream_name" description:"流水线名"`
	StreamDesc string    `json:"stream_desc" gorm:"column:stream_desc" description:"描述"`
	UpdatedAt  time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	CreatedAt  time.Time `json:"create_at" gorm:"column:create_at" description:"创建时间"`
	IsDelete   int8      `json:"is_delete" gorm:"column:is_delete" description:"是否删除;0:否;1:是"`
}

// TableName ...
func (t *StreamInfo) TableName() string {
	return "chain_stream"
}

// Find 根据ID查找
func (t *StreamInfo) Find(c *gin.Context, tx *gorm.DB, search *StreamInfo) (*StreamInfo, error) {
	model := &StreamInfo{}
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(model).Error
	return model, err
}

// Save 保存对象t到DB
func (t *StreamInfo) Save(c *gin.Context, tx *gorm.DB) error {
	if err := tx.SetCtx(public.GetGinTraceContext(c)).Save(t).Error; err != nil {
		return err
	}
	return nil
}

// PageList 分页查询=>得服务数组
func (t *StreamInfo) PageList(c *gin.Context, tx *gorm.DB, search *dto.ServiceListInput) ([]StreamInfo, int64, error) {
	var total int64 = 0
	list := []StreamInfo{}
	offset := (search.PageNo - 1) * search.PageSize //第一页的话就不用偏移了,直接limit查
	query := tx.SetCtx(public.GetGinTraceContext(c))
	query = query.Table(t.TableName()).Where("is_delete=0")
	if search.Info != "" { ///模糊查询
		query = query.Where("stream_name like ?", "%"+search.Info+"%")
	}
	if err := query.Limit(search.PageSize).Offset(offset).Order("id desc").Find(&list).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	query.Limit(search.PageSize).Offset(offset).Count(&total)
	return list, total, nil
}
