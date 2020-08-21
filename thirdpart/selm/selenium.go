package selm

import (
	"fmt"
	"time"

	"github.com/jstang9527/gateway/dto"

	"github.com/jstang9527/gateway/dao"

	"github.com/tebeka/selenium"
)

// BlockTask 组件对象
type BlockTask struct {
	BlockInfo *dao.BlockInfo
	Actions   []*dao.ActionItem
	Params    *dto.WebhookInput
}

// CreateTask 1.创建测试任务
func CreateTask(detailList []*dao.BlockDetail, params *dto.WebhookInput) (err error) {
	caps.AddChrome(chromeCaps)
	wd, err := selenium.NewRemote(caps, "http://127.0.0.1:9515/wd/hub")
	if err != nil {
		fmt.Printf("unable create browser, err: %v", err)
		return
	}
	wd.SetImplicitWaitTimeout(time.Second * 10)
	go run(detailList, params, wd)
	return
}

// 遍历多条链 2
func run(detailList []*dao.BlockDetail, params *dto.WebhookInput, wd selenium.WebDriver) {
	defer wd.Quit()
	for _, block := range detailList {
		// 一条一条链执行
		blockTask := NewBlockTask(block.Info, block.Actions, params)
		blockTask.ExecBlock(wd)

	}
}

// NewBlockTask 3.创建组件对象
func NewBlockTask(blockInfo *dao.BlockInfo, actions []*dao.ActionItem, params *dto.WebhookInput) *BlockTask {
	return &BlockTask{BlockInfo: blockInfo, Actions: actions, Params: params}
}

// ExecBlock 4.执行单条链的所有动作
func (bt *BlockTask) ExecBlock(wd selenium.WebDriver) {
	var actionTask = &ActionTask{}
	for _, action := range bt.Actions {
		// 一个动作一个动作执行
		actionTask = NewActionTask(bt, wd)
		if err := actionTask.ExecAction(action); err != nil {
			// 一条出错，终止整条链
			break
		}
		time.Sleep(time.Second)
	}
	aa := &dao.ActionItem{
		EventType: 1,
		URL:       "/#/overview",
	}
	actionTask.EntryAction(aa)
	wd.Refresh()
}
