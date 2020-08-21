package dao

import (
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
)

// BlockDetail 多表查询后形成的服务详细信息结构体
type BlockDetail struct {
	Info    *BlockInfo    `json:"info" description:"基本信息"`
	Actions []*ActionItem `json:"actions" description:"动作列表"`
}

// GetBlockDetail 组件详情 ==> 获得多表查询后的组成的结构体blockDetail
// 通过blockInfo的外键ID查其他表的关联信息
// 得到的结果封装到完全体BlockDetail
func (s *BlockInfo) GetBlockDetail(c *gin.Context, tx *gorm.DB, search *BlockInfo) (*BlockDetail, error) {
	if search.BlockName == "" {
		info, err := s.Find(c, tx, search)
		if err != nil {
			return nil, err
		}
		search = info
	}

	action := &ActionItem{BlockID: search.ID}
	actionList, err := action.GetActionsByBlockID(c, tx, action)
	if err != nil && err != gorm.ErrRecordNotFound { //找不到也是正确的
		return nil, err
	}

	detail := &BlockDetail{Info: search, Actions: actionList}
	return detail, nil
}
