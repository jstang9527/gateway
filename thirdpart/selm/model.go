package selm

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/jstang9527/gateway/dto"

	"github.com/jstang9527/gateway/thirdpart/es"

	"github.com/jstang9527/gateway/public"
	"github.com/jstang9527/opentest/utils/tools"
	"github.com/tebeka/selenium"

	"github.com/jstang9527/gateway/dao"
)

// ActionTask ...
type ActionTask struct {
	ProjectID   string        //md5
	ProjectName string        //detux
	ProjectAddr string        //https://gitlab.com/xxx
	WebAddr     string        // http://172.31.50.252:65000
	FuncName    string        //创建蜜罐
	Message     string        // OpenHomePage, err:nil, output: null  *
	Screenshot  string        //https://.....                         *
	Priority    int           // High
	Status      int           // Pass   123                           *
	Duration    time.Duration //                                      *
	Wd          selenium.WebDriver
}

// NewActionTask ...
func NewActionTask(bt *BlockTask, wd selenium.WebDriver) *ActionTask {
	return &ActionTask{
		ProjectID:   public.MD5(fmt.Sprintf("%v%v%v%v", bt.Params.ProjectName, bt.Params.ProjectAddr, bt.Params.WebAddr, bt.Params.StreamID)),
		ProjectName: bt.Params.ProjectName,
		ProjectAddr: bt.Params.ProjectAddr,
		WebAddr:     bt.Params.WebAddr,
		FuncName:    bt.BlockInfo.BlockName,
		Priority:    bt.BlockInfo.Priority,
		Wd:          wd,
	}
}

// ExecAction 执行具体动作
func (at *ActionTask) ExecAction(action *dao.ActionItem) (err error) {
	start := time.Now()
	status := 1
	message := "" //接收抓取的数据

	//1.根据操作类型执行对应操作函数, 输入事件|点击事件|页面跳转事件|指标监控事件|数据抓取事件
	switch action.EventType {
	case entry:
		err = at.EntryAction(action)
	case input:
		err = at.InputAction(action)
	case click:
		err = at.ClickAction(action)
	case monitor:
		err = at.MonitorAction(action)
	case parse:
		message, err = at.ParseAction(action)
	default:
		err = fmt.Errorf("params eleType not in [entry,input,click,monitor,parse]")
		status = 2
	}
	//2.判断出错与否
	if err != nil { //有错
		if status != 2 { //非人为写错
			status = 3
		}
	}
	//3.把正确或者错误的结果存起来
	if message != "" { //output输出
		message = fmt.Sprintf("%v, quota: %v, err: %v, output: %v.", action.ActionName, action.ElementValue, err, message)
	} else {
		message = fmt.Sprintf("%v, quota: %v, err: %v", action.ActionName, action.ElementValue, err)
	}
	screenshot := public.MD5(message) //截图名

	// 保存截图、异步推ES
	at.Save(message, screenshot, status, time.Now().Sub(start))
	return
}

// Save ...
func (at *ActionTask) Save(message, screenshot string, status int, dura time.Duration) {
	at.Message = message
	at.Screenshot = screenshot + ".png"
	at.Status = status
	at.Duration = dura
	fmt.Printf("%#v\n", at)
	//同步保存截图，否则at.wd会报错
	at.saveScreenshot()
	//异步推ES
	ti := dto.NewTestItem(at.ProjectID, at.ProjectName, at.ProjectAddr, at.WebAddr, at.FuncName, at.Message, at.Screenshot, at.Priority, at.Status, at.Duration)
	go es.SendToESChan(ti)
}

// 保存截图 ...(保存路径为服务器内文件路径, web访问路径非此路径)
func (at *ActionTask) saveScreenshot() {
	//创建存储项目id的截图目录
	dirPath := "./resources/" + at.ProjectID + "/"
	if bl, _ := tools.PathExists(dirPath); !bl {
		os.Mkdir(dirPath, 0666)
	}
	imgpath := dirPath + at.Screenshot
	time.Sleep(time.Millisecond * 200) //等0.8秒再截图

	if b, err := at.Wd.Screenshot(); err != nil {
		fmt.Printf("failed output screenshop by selenium, err: %v\n", err)
	} else {
		if err := ioutil.WriteFile(imgpath, b, 0666); err != nil {
			fmt.Printf("cann't save screenshot in local, err: %v\n", err)
			return
		}
	}
}

// EntryAction ...
func (at *ActionTask) EntryAction(action *dao.ActionItem) (err error) {
	url := at.WebAddr + action.URL // http://172.31.50.39:65000
	return at.Wd.Get(url)
}

// InputAction ...
func (at *ActionTask) InputAction(action *dao.ActionItem) (err error) {
	var ele selenium.WebElement
	//1.先判断是id还是xpath
	if action.SearchType == 1 {
		ele, err = at.Wd.FindElement(selenium.ByID, action.ElementID)
	} else if action.SearchType == 2 {
		ele, err = at.Wd.FindElement(selenium.ByXPATH, action.XPath)
	} else {
		return fmt.Errorf("unable found expect search type: %v; [1:byid, 2:xpath]", action.SearchType)
	}
	if err != nil {
		return fmt.Errorf("element not found, info: %v", err)
	}
	return ele.SendKeys(action.ElementValue)
}

// ClickAction ...
func (at *ActionTask) ClickAction(action *dao.ActionItem) (err error) {
	var ele selenium.WebElement
	//1.先判断是id还是xpath
	if action.SearchType == 1 {
		ele, err = at.Wd.FindElement(selenium.ByID, action.ElementID)
	} else if action.SearchType == 2 {
		ele, err = at.Wd.FindElement(selenium.ByXPATH, action.XPath)
	} else {
		return fmt.Errorf("unable found expect search type: %v; [1:byid, 2:xpatj]", action.SearchType)
	}
	if err != nil {
		return fmt.Errorf("element not found, info: %v", err)
	}
	return ele.Click()
}

// MonitorAction ...
func (at *ActionTask) MonitorAction(action *dao.ActionItem) (err error) {
	var ele selenium.WebElement
	var quota string

	// 1.设置元素检索超时时间
	if action.Timeout > 0 {
		at.Wd.SetImplicitWaitTimeout(time.Second * time.Duration(action.Timeout))
	}
	defer func() { at.Wd.SetImplicitWaitTimeout(time.Second * 10) }()

	// 2.开始检索数据
	// 2.1.先判断是id还是xpath
	if action.SearchType == 1 {
		ele, err = at.Wd.FindElement(selenium.ByID, action.ElementID)
	} else if action.SearchType == 2 {
		ele, err = at.Wd.FindElement(selenium.ByXPATH, action.XPath)
	} else {
		return fmt.Errorf("unable found expect search type: %v; [1:byid, 2:xpatj]", action.SearchType)
	}
	// 2.2 元素找不到
	if err != nil {
		if action.AllowErr == 1 { //允许
			return nil
		}
		return fmt.Errorf("element not found, info: %v", err)
	}
	// 2.3 元素找到了
	quota, err = ele.Text()
	if err != nil {
		return fmt.Errorf("unable get text of element,info: %v", err)
	}
	reg := regexp.MustCompile(`\s{2,}|\n|\t|\r`)
	quota = reg.ReplaceAllString(quota, "")
	quota = strings.TrimSpace(quota)
	quota = strings.ToLower(quota)
	// 4.是否满足指标
	if quota == action.ElementValue {
		fmt.Printf("excpet: %v, output: %v\n", action.ElementValue, quota)
		return nil
	}
	// 5.不满足
	return fmt.Errorf("excpet: %v, output: %v", action.ElementValue, quota)
}

// ParseAction 抓数据, 看用户的value,没设置的话空则info,有设置的话Warning
func (at *ActionTask) ParseAction(action *dao.ActionItem) (result string, err error) {
	var ele selenium.WebElement
	//1.先判断是id还是xpath
	if action.SearchType == 1 {
		ele, err = at.Wd.FindElement(selenium.ByID, action.ElementID)
	} else if action.SearchType == 2 {
		ele, err = at.Wd.FindElement(selenium.ByXPATH, action.XPath)
	} else {
		return "", fmt.Errorf("unable found expect search type: %v; [1:byid, 2:xpatj]", action.SearchType)
	}
	if err != nil {
		return "", fmt.Errorf("element not found, info: %v", err)
	}

	//2.抓数据
	rows, err := ele.FindElements(selenium.ByTagName, "tr")
	if err != nil {
		err = fmt.Errorf("element not found")
		return
	}
	for _, row := range rows {
		items, _ := row.FindElements(selenium.ByTagName, "td")
		for _, item := range items {
			text, _ := item.Text()
			result = fmt.Sprintf("%v, %v", result, strings.TrimSpace(text))
		}
	}
	reg := regexp.MustCompile(`\s{2,}|\n|\t|\r`)
	result = reg.ReplaceAllString(result, "")
	return
}
