package dao

import (
	"time"

	"github.com/jstang9527/gateway/dto"
	"github.com/jstang9527/gateway/public"

	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
)

// BlockInfo 组件由多个动作组成
type BlockInfo struct {
	ID        int64     `json:"id" gorm:"primary_key"`
	BlockName string    `json:"block_name" gorm:"column:block_name" description:"组件名"`
	Priority  int       `json:"priority" gorm:"column:priority" description:"重要性"`
	Expect    string    `json:"expect" gorm:"column:expect" description:"期望值"`
	UpdatedAt time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	CreatedAt time.Time `json:"create_at" gorm:"column:create_at" description:"创建时间"`
	IsDelete  int8      `json:"is_delete" gorm:"column:is_delete" description:"是否删除;0:否;1:是"`
}

// TableName ...
func (t *BlockInfo) TableName() string {
	return "chain_block"
}

// Find 根据ID查找
func (t *BlockInfo) Find(c *gin.Context, tx *gorm.DB, search *BlockInfo) (*BlockInfo, error) {
	model := &BlockInfo{}
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(model).Error
	return model, err
}

// Save 保存对象t到DB
func (t *BlockInfo) Save(c *gin.Context, tx *gorm.DB) error {
	if err := tx.SetCtx(public.GetGinTraceContext(c)).Save(t).Error; err != nil {
		return err
	}
	return nil
}

// PageList 分页查询=>得服务数组
func (t *BlockInfo) PageList(c *gin.Context, tx *gorm.DB, search *dto.ServiceListInput) ([]BlockInfo, int64, error) {
	var total int64 = 0
	list := []BlockInfo{}
	offset := (search.PageNo - 1) * search.PageSize //第一页的话就不用偏移了,直接limit查
	query := tx.SetCtx(public.GetGinTraceContext(c))
	query = query.Table(t.TableName()).Where("is_delete=0")
	if search.Info != "" { ///模糊查询
		query = query.Where("block_name like ?", "%"+search.Info+"%")
	}
	if err := query.Limit(search.PageSize).Offset(offset).Order("id desc").Find(&list).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	query.Limit(search.PageSize).Offset(offset).Count(&total)
	return list, total, nil
}
