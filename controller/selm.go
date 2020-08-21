package controller

import (
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"github.com/jstang9527/gateway/dao"
	"github.com/jstang9527/gateway/dto"
	"github.com/jstang9527/gateway/middleware"
	"github.com/jstang9527/gateway/thirdpart/selm"
)

// SelmController 主机控制器结果体
type SelmController struct{}

// SelmRegister 主机控制器
func SelmRegister(group *gin.RouterGroup) {
	SelmCtl := &SelmController{}
	group.POST("/task", SelmCtl.ExecuteTestTask)
}

// ExecuteTestTask godoc
// @Summary webhook
// @Description webhook
// @Tags Web自动化测试
// @ID /selm/task/post
// @Accept json
// @Produce json
// @Param project_name query string true "项目名"
// @Param project_addr query string false "项目地址"
// @Param web_addr query string true "web地址"
// @Param sync_num query int false "并发数"
// @Param search_timeout query int false "检索元素超时时间"
// @Param stream_id query int true "流水线ID"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /selm/task [post]
func (s *SelmController) ExecuteTestTask(c *gin.Context) {
	inputParams := &dto.WebhookInput{}
	if err := inputParams.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	//1.根据stream_id从多对多表中取组件ID列表
	bsmi := dao.BlockStreamMultiInfo{StreamID: inputParams.StreamID}
	list, err := bsmi.Find(c, tx, &bsmi)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	//2.遍历list的blockid
	var detailList []*dao.BlockDetail
	for _, item := range *list {
		//2.1 根据block_id,查block基本信息
		blockInfo := &dao.BlockInfo{ID: item.BlockID}
		blockInfo, err = blockInfo.Find(c, tx, blockInfo)
		if err != nil {
			middleware.ResponseError(c, 2003, err)
			return
		}
		//2.2 sdetail包含组件自身信息及其所有动作信息
		sdetail, err := blockInfo.GetBlockDetail(c, tx, blockInfo)
		if err != nil {
			middleware.ResponseError(c, 2004, err)
			return
		}
		detailList = append(detailList, sdetail)
	}
	// fmt.Printf("%#v,%#v", detailList[0].Info, detailList[1].Actions)

	// CreateTask是个异步，会立即返回环境初始化是否正常
	if err = selm.CreateTask(detailList, inputParams); err != nil {
		middleware.ResponseError(c, 2005, err)
		return
	}

	out := "success"
	middleware.ResponseSuccess(c, out)
}
