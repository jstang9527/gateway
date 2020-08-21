package controller

import (
	"fmt"

	"github.com/e421083458/golang_common/lib"
	"github.com/jstang9527/gateway/dao"

	"github.com/gin-gonic/gin"
	"github.com/jstang9527/gateway/dto"
	"github.com/jstang9527/gateway/middleware"
)

// ChainController ...
type ChainController struct{}

// ChainRegister ...
func ChainRegister(group *gin.RouterGroup) {
	admin := &ChainController{}
	group.POST("/block", admin.ChainAddBlock)
	group.GET("/block", admin.BlockDetail)
	group.GET("/block_list", admin.BlockList)
	group.POST("/stream", admin.ChainAddStream)
	group.GET("/stream_list", admin.StreamList)
	group.PUT("/block", admin.ChainUpdateBlock)
}

// ChainAddBlock godoc
// @Summary 添加组件
// @Description 添加组件
// @Tags Web自动化测试
// @ID /chain/block/post
// @Accept json
// @Produce json
// @Param body body dto.ChainAddBlockInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /chain/block [post]
func (admin *ChainController) ChainAddBlock(c *gin.Context) {
	//1. 请求参数初步校验(必填)
	inputParams := &dto.ChainAddBlockInput{}
	if err := inputParams.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	// fmt.Printf("%#v\n", inputParams)
	// fmt.Printf("%#v\n", inputParams.ActionChain[0])
	// fmt.Printf("%#v\n", inputParams.ActionChain[1])
	//---------- 要不要进行信息处理？
	//2. 直接存数据库
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	tx = tx.Begin() //开始事务
	//2.1 组件信息入库
	blockModel := &dao.BlockInfo{
		BlockName: inputParams.BlockName,
		Expect:    inputParams.Expect,
		Priority:  inputParams.Priority,
	}
	if err := blockModel.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2002, err)
		return
	}
	//2.2 动作信息入库
	for i := 0; i < len(inputParams.ActionChain); i++ {
		fmt.Println(i)
		item := inputParams.ActionChain[i]
		actionModel := &dao.ActionItem{
			BlockID:      blockModel.ID,
			ActionName:   item.ActionName,
			AllowErr:     item.AllErr,
			ElementID:    item.ElementID,
			ElementValue: item.ElementValue,
			EventType:    item.EventType,
			SearchType:   item.SearchType,
			Timeout:      item.Timeout,
			Timestamp:    item.Timestamp,
			URL:          item.URL,
			XPath:        item.XPath,
		}
		if err := actionModel.Save(c, tx); err != nil {
			tx.Rollback()
			middleware.ResponseError(c, 2003, err)
			return
		}
	}
	tx.Commit() //提交事务
	//6. 返回信息
	out := fmt.Sprintf("Create ActionBlock [%v] Success.", blockModel.BlockName)
	middleware.ResponseSuccess(c, out)
}

// BlockList godoc
// @Summary 组件列表
// @Description 组件列表
// @Tags Web自动化测试
// @ID /chain/block_list
// @Accept json
// @Produce json
// @Param info query string false "关键词"
// @Param page_size query int true "每页个数"
// @Param page_no query int true "当前页数"
// @Success 200 {object} middleware.Response{data=dto.BlockListOutput} "success"
// @Router /chain/block_list [get]
func (admin *ChainController) BlockList(c *gin.Context) {
	inputParams := &dto.ServiceListInput{}
	if err := inputParams.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	//从db中分页读取基本信息
	blockInfo := &dao.BlockInfo{}
	blockList, total, err := blockInfo.PageList(c, tx, inputParams)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	//格式化输出信息
	outList := []dto.BlockListItemOutput{} //这个结构体是面向前端接口的
	for _, item := range blockList {
		outItem := dto.BlockListItemOutput{
			ID:        item.ID,
			BlockName: item.BlockName,
			Priority:  item.Priority,
			Expect:    item.Expect,
		}
		outList = append(outList, outItem)
	}
	out := &dto.BlockListOutput{Total: total, List: outList}
	middleware.ResponseSuccess(c, out)
}

// BlockDetail godoc
// @Summary 组件详情
// @Description 组件详情
// @Tags Web自动化测试
// @ID /chain/block/get
// @Accept json
// @Produce json
// @Param id query string true "组件ID"
// @Success 200 {object} middleware.Response{data=dao.BlockDetail} "success"
// @Router /chain/block [get]
func (admin *ChainController) BlockDetail(c *gin.Context) {
	params := &dto.ServiceDeleteInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	//从db中读取组件信息
	blockInfo := &dao.BlockInfo{ID: params.ID}
	blockInfo, err = blockInfo.Find(c, tx, blockInfo)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	//从DB在读取组件、及其所有动作信息
	sdetail, err := blockInfo.GetBlockDetail(c, tx, blockInfo)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	out := sdetail
	middleware.ResponseSuccess(c, out)
}

// ChainAddStream godoc
// @Summary 添加流水线
// @Description 添加流水线
// @Tags Web自动化测试
// @ID /chain/stream/post
// @Accept json
// @Produce json
// @Param body body dto.ChainAddStreamInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /chain/stream [post]
func (admin *ChainController) ChainAddStream(c *gin.Context) {
	//1. 请求参数初步校验(必填)
	inputParams := &dto.ChainAddStreamInput{}
	if err := inputParams.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	//2. 直接存数据库
	//2.0 连接池
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	tx = tx.Begin() //开始事务
	//2.1 流水线信息入库
	StreamModel := &dao.StreamInfo{
		StreamName: inputParams.StreamName,
		StreamDesc: inputParams.StreamDesc,
	}
	if err := StreamModel.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2002, err)
		return
	}
	//2.2 多对多信息入库
	for i := 0; i < len(inputParams.BlockList); i++ {
		blockID := inputParams.BlockList[i]
		multiModel := &dao.BlockStreamMultiInfo{
			BlockID:  blockID,
			StreamID: StreamModel.ID,
		}
		if err := multiModel.Save(c, tx); err != nil {
			tx.Rollback()
			middleware.ResponseError(c, 2003, err)
			return
		}
	}
	tx.Commit() //提交事务
	//6. 返回信息
	out := fmt.Sprintf("Create ChainStream [%v] Success.", StreamModel.StreamName)
	middleware.ResponseSuccess(c, out)
}

// StreamList godoc
// @Summary 流水线列表
// @Description 流水线列表
// @Tags Web自动化测试
// @ID /chain/stream_list
// @Accept json
// @Produce json
// @Param info query string false "关键词"
// @Param page_size query int true "每页个数"
// @Param page_no query int true "当前页数"
// @Success 200 {object} middleware.Response{data=dto.StreamListOutput} "success"
// @Router /chain/stream_list [get]
func (admin *ChainController) StreamList(c *gin.Context) {
	inputParams := &dto.ServiceListInput{}
	if err := inputParams.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	//从db中分页读取基本信息
	streamInfo := &dao.StreamInfo{}
	blockList, total, err := streamInfo.PageList(c, tx, inputParams)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	//格式化输出信息
	outList := []dto.StreamListItemOutput{} //这个结构体是面向前端接口的
	for _, item := range blockList {
		outItem := dto.StreamListItemOutput{
			ID:         item.ID,
			StreamName: item.StreamName,
			StreamDesc: item.StreamDesc,
			API:        fmt.Sprintf("http://172.31.50.254:8880/selm/task?stream_id=%v&web_addr=http://ip:port&project_name=example", item.ID),
		}
		outList = append(outList, outItem)
	}
	out := &dto.StreamListOutput{Total: total, List: outList}
	middleware.ResponseSuccess(c, out)
}

// ChainUpdateBlock godoc
// @Summary 组件更新
// @Description 组件更新
// @Tags Web自动化测试
// @ID /chain/block/update
// @Accept  json
// @Produce  json
// @Param body body dto.ChainUpdateBlockInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /chain/block [put]
func (admin *ChainController) ChainUpdateBlock(c *gin.Context) {
	params := &dto.ChainUpdateBlockInput{}
	if err := params.GetValidParams(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	// 1.拿组件ID去查表
	blockInfo := &dao.BlockInfo{ID: params.ID}
	// 2.根据组件ID查询其动作列表对象
	blockDetail, err := blockInfo.GetBlockDetail(c, tx, blockInfo)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	blockInfo = blockDetail.Info
	actions := blockDetail.Actions //数据库中的旧动作列表

	// fmt.Printf("%#v\n", blockInfo)
	// fmt.Printf("%#v\n", actions[1])
	// 3.开始修改、保存对象
	tx = tx.Begin()
	// //3.1 保存组件自身信息
	blockInfo.BlockName = params.BlockName
	blockInfo.Expect = params.Expect
	blockInfo.Priority = params.Priority
	if err := blockInfo.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2003, err)
		return
	}
	// 若是修改action，action_id必不为空
	InputActions := params.ActionChain
	for _, item := range InputActions {
		if item.ID > 0 { //更新操作
			for i, oldaction := range actions {
				if oldaction.ID == item.ID {
					oldaction.ActionName = item.ActionName
					oldaction.AllowErr = item.AllErr
					oldaction.ElementID = item.ElementID
					oldaction.ElementValue = item.ElementValue
					oldaction.EventType = item.EventType
					oldaction.SearchType = item.SearchType
					oldaction.Timeout = item.Timeout
					oldaction.Timestamp = item.Timestamp
					oldaction.URL = item.URL
					oldaction.XPath = item.XPath
					if err := oldaction.Save(c, tx); err != nil {
						tx.Rollback()
						middleware.ResponseError(c, 2004, err)
						return
					}
					//修改过的就去掉
					actions = append(actions[:i], actions[i+1:]...)
					break
				}
			}

		} else { //新增操作
			newItem := &dao.ActionItem{
				BlockID:      blockInfo.ID,
				ActionName:   item.ActionName,
				AllowErr:     item.AllErr,
				ElementID:    item.ElementID,
				ElementValue: item.ElementValue,
				EventType:    item.EventType,
				SearchType:   item.SearchType,
				Timeout:      item.Timeout,
				Timestamp:    item.Timestamp,
				URL:          item.URL,
				XPath:        item.XPath,
			}
			if err := newItem.Save(c, tx); err != nil {
				tx.Rollback()
				middleware.ResponseError(c, 2005, err)
				return
			}
		}

	}

	// 剩下的actions都是要被删除的
	for _, action := range actions {
		if err := action.Delete(c, tx); err != nil {
			tx.Rollback()
			middleware.ResponseError(c, 2005, err)
			return
		}
	}

	tx.Commit()
	middleware.ResponseSuccess(c, "Update Success.")
	return
}
