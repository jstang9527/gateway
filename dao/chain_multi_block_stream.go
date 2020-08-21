package dao

import (
	"github.com/jstang9527/gateway/public"

	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
)

// BlockStreamMultiInfo 组件由多个动作组成
type BlockStreamMultiInfo struct {
	ID       int64 `json:"id" gorm:"primary_key"`
	BlockID  int64 `json:"block_id" gorm:"column:block_id" description:"组件ID"`
	StreamID int64 `json:"stream_id" gorm:"column:stream_id" description:"流水线ID"`
}

// TableName ...
func (t *BlockStreamMultiInfo) TableName() string {
	return "chain_multi_block_stream"
}

// Find 根据ID查找
func (t *BlockStreamMultiInfo) Find(c *gin.Context, tx *gorm.DB, search *BlockStreamMultiInfo) (*[]BlockStreamMultiInfo, error) {
	var models = []BlockStreamMultiInfo{}
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(&models).Error
	return &models, err
}

// Save 保存对象t到DB
func (t *BlockStreamMultiInfo) Save(c *gin.Context, tx *gorm.DB) error {
	if err := tx.SetCtx(public.GetGinTraceContext(c)).Save(t).Error; err != nil {
		return err
	}
	return nil
}
