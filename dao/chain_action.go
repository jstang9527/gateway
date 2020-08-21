package dao

import (
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"github.com/jstang9527/gateway/public"
)

// ActionItem 单个动作
type ActionItem struct {
	ID           int64  `json:"id" gorm:"primary_key"`
	BlockID      int64  `json:"block_id" gorm:"column:block_id" description:"服务id"`
	ActionName   string `json:"action_name" gorm:"column:action_name" description:"组件名"`
	AllowErr     int    `json:"allow_err" gorm:"column:allow_err" description:"错误为真"`
	ElementID    string `json:"element_id" gorm:"column:element_id" description:"元素ID"`
	ElementValue string `json:"element_value" gorm:"column:element_value" description:"元素值"`
	EventType    int    `json:"event_type" gorm:"column:event_type" description:"事件类型"`
	SearchType   int    `json:"search_type" gorm:"column:search_type" description:"检索方式"`
	Timeout      int    `json:"timeout" gorm:"column:timeout" description:"超时时间"`
	Timestamp    string `json:"timestamp" gorm:"column:timestamp" description:"中文时间"`
	URL          string `json:"url" gorm:"column:url" description:"资源定位符"`
	XPath        string `json:"xpath" gorm:"column:xpath" description:"XPath"`
}

// TableName ...
func (t *ActionItem) TableName() string {
	return "chain_action"
}

// Find ...
func (t *ActionItem) Find(c *gin.Context, tx *gorm.DB, search *ActionItem) (*ActionItem, error) {
	model := &ActionItem{}
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(model).Error
	return model, err
}

// Save ...
func (t *ActionItem) Save(c *gin.Context, tx *gorm.DB) error {
	if err := tx.SetCtx(public.GetGinTraceContext(c)).Save(t).Error; err != nil {
		return err
	}
	return nil
}

// Delete 删除单个动作action
func (t *ActionItem) Delete(c *gin.Context, tx *gorm.DB) error {
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where("id=?", t.ID).Delete(t).Error
	return err
}

// GetActionsByBlockID 获取组件下的所有Action动作
func (t *ActionItem) GetActionsByBlockID(c *gin.Context, tx *gorm.DB, search *ActionItem) ([]*ActionItem, error) {
	list := []*ActionItem{}
	query := tx.SetCtx(public.GetGinTraceContext(c))
	query = query.Table(t.TableName()).Where("block_id = ?", search.BlockID)

	// if err := query.Order("id desc").Find(&list).Error; err != nil && err != gorm.ErrRecordNotFound {
	if err := query.Find(&list).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return list, nil
}
