package dto

import (
	"github.com/gin-gonic/gin"
	"github.com/jstang9527/gateway/public"
)

// ActionItemInput 单个动作
type ActionItemInput struct {
	ID           int64  `json:"id" form:"id" comment:"新建为空，修改可空" example:"" validate:""`
	ActionName   string `json:"action_name" form:"action_name" comment:"动作名" example:"" validate:"required"`            //动作名
	AllErr       int    `json:"all_err" form:"all_err" comment:"错误为真:1激活,0默认" example:"" validate:""`                   //错误为真
	ElementID    string `json:"element_id" form:"element_id" comment:"元素ID" example:"" validate:""`                     //元素ID
	ElementValue string `json:"element_value" form:"element_value" comment:"元素值" example:"" validate:""`                //元素值
	EventType    int    `json:"event_type" form:"event_type" comment:"事件类型:1-页面,2-输入,3-点击,4-抓取" example:"" validate:""` //事件类型
	SearchType   int    `json:"search_type" form:"search_type" comment:"元素检索方式:xpath|byid" example:"" validate:""`      //元素检索方式
	Timeout      int    `json:"timeout" form:"timeout" comment:"检索超时时间,0值设置系统默认值" example:"" validate:""`               //检索超时时间
	Timestamp    string `json:"timestamp" form:"timestamp" comment:"时间戳" example:"" validate:""`                        //时间戳
	URL          string `json:"url" form:"url" comment:"目标网址" example:"" validate:""`                                   //目标网址
	XPath        string `json:"xpath" form:"xpath" comment:"XPath" example:"" validate:""`                              //XPath
}

// ChainAddBlockInput 由多个动作组成的动作链(组件)
type ChainAddBlockInput struct {
	BlockName   string             `json:"block_name" form:"block_name" comment:"组件名" example:"" validate:"required"`          // 组件名
	Priority    int                `json:"priority" form:"priority" comment:"重要性:1-Low,2-Nomal,3-High" example:"" validate:""` // 重要性
	Expect      string             `json:"expect" form:"expect" comment:"期望值" example:"" validate:""`                          // 期望值
	ActionChain []*ActionItemInput `json:"action_chain" form:"action_chain" comment:"动作列表" example:"" validate:""`             // 动作列表
}

// ChainAddStreamInput ...
type ChainAddStreamInput struct {
	StreamName string  `json:"stream_name" form:"stream_name" comment:"流水线名" example:"" validate:"required"` // 流水线名
	StreamDesc string  `json:"stream_desc" form:"stream_desc" comment:"描述" example:"" validate:""`           // 描述
	BlockList  []int64 `json:"block_list" form:"block_list" comment:"组件列表" example:"" validate:""`           // 组件列表
}

// BindValidParam 校验新增参数,绑定结构体,校验参数
func (s *ChainAddStreamInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, s)
}

// BindValidParam 校验新增参数,绑定结构体,校验参数
func (s *ChainAddBlockInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, s)
}

// BindValidParam 校验新增参数,绑定结构体,校验参数
func (s *ActionItemInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, s)
}

// BlockListOutput ...
type BlockListOutput struct {
	Total int64                 `json:"total" form:"total" comment:"总数" example:"" validate:""` //总数
	List  []BlockListItemOutput `json:"list" form:"list" comment:"列表" example:"" validate:""`   //列表
}

// BlockListItemOutput 对象item信息
type BlockListItemOutput struct {
	ID        int64  `json:"id" form:"id"`                 //id
	BlockName string `json:"block_name" form:"block_name"` //组件名称
	Priority  int    `json:"priority" form:"priority"`     //重要性
	Expect    string `json:"expect" form:"expect"`         //期望值
}

// StreamListOutput ...
type StreamListOutput struct {
	Total int64                  `json:"total" form:"total" comment:"总数" example:"" validate:""` //总数
	List  []StreamListItemOutput `json:"list" form:"list" comment:"列表" example:"" validate:""`   //列表
}

// StreamListItemOutput 对象item信息
type StreamListItemOutput struct {
	ID         int64  `json:"id" form:"id"`                   //id
	StreamName string `json:"stream_name" form:"stream_name"` //流水线名称
	StreamDesc string `json:"stream_desc" form:"stream_desc"` //描述
	API        string `json:"api" form:"api"`                 //调用API,是给Jenkins调用滴
}

// ChainUpdateBlockInput 除了组件ID,其他均可更改
type ChainUpdateBlockInput struct {
	ID          int64              `json:"id" form:"id" comment:"组件ID" validate:"required"`
	BlockName   string             `json:"block_name" form:"block_name" comment:"组件名" example:"" validate:"required"`          // 组件名
	Priority    int                `json:"priority" form:"priority" comment:"重要性:1-Low,2-Nomal,3-High" example:"" validate:""` // 重要性
	Expect      string             `json:"expect" form:"expect" comment:"期望值" example:"" validate:""`                          // 期望值
	ActionChain []*ActionItemInput `json:"action_chain" form:"action_chain" comment:"动作列表" example:"" validate:""`             // 动作列表
}

// GetValidParams ...
func (params *ChainUpdateBlockInput) GetValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}
