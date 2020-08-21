package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/TruthHun/html2json/html2json"
	"github.com/jstang9527/gateway/thirdpart/selm2"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

const (
	//设置常量 分别设置chromedriver.exe的地址和本地调用端口
	seleniumPath = `/root/chromedriver`
	port         = 9515
)

var (
	chromeCaps = chrome.Capabilities{
		// Prefs: imgCaps,
		Path: "",
		Args: []string{
			// "--headless",  //不开启浏览器
			"--start-maximized",
			"--window-size=1200x600",
			"--no-sandbox", //非root可运行
			"--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36",
			"--disable-gpu",
			"--disable-impl-side-painting",
			"--disable-gpu-sandbox",
			"--disable-accelerated-2d-canvas",
			"--disable-accelerated-jpeg-decoding",
			"--test-type=ui",
			"--ignore-certificate-errors",
		},
	}
)

// Init ...
func Init() (err error) {
	ops := []selenium.ServiceOption{}
	service, err = selenium.NewChromeDriverService(seleniumPath, port, ops...)
	if err != nil {
		return
	}

	caps := selenium.Capabilities{"browserName": "chrome"}
	caps.AddChrome(chromeCaps)
	wd, err = selenium.NewRemote(caps, "http://127.0.0.1:9515/wd/hub")
	return
}

var wd selenium.WebDriver
var service *selenium.Service

func main() {
	// --------------------------init------------------------------------
	if err := Init(); err != nil {
		fmt.Println("Failed init selenium service, info: ", err)
		return
	}
	defer wd.Quit()
	defer service.Stop()
	// --------------------------启动webdriver---------------------------
	if err := wd.Get("https://172.31.50.39:65000/#/config-manage/sys"); err != nil {
		fmt.Println("Failed to open index page,", err)
		return
	}
	// --------------------------run------------------------------------
	run()
	// --------------------------结束所有--------------------------------
	time.Sleep(time.Second * 10)

}

func run() (err error) {
	time.Sleep(time.Second * 3)
	html, err := wd.PageSource()
	if err != nil {
		return
	}
	// 1、匹配精简目标
	flysnowRegexp := regexp.MustCompile(`<div id="root">(.*)</div>`)
	params := flysnowRegexp.FindStringSubmatch(html)
	// 2、解析HTML成对象
	html = strings.Join(params, "")
	reg := regexp.MustCompile(`\s{2,}|\n|\t|\r`)
	html = reg.ReplaceAllString(html, "")
	appTags := html2json.GetTags(html2json.TagUniAPP)
	rt := html2json.New(appTags)
	nodes, err := rt.Parse(html, "")
	if err != nil {
		fmt.Println("Failed Parse Html Nodes,", err)
		return
	}
	// 3、转json,再转自家结构体
	jsonstr := selm2.ToJSON(nodes)
	objs := []selm2.H2j{}
	if err = json.Unmarshal([]byte(jsonstr), &objs); err != nil {
		fmt.Println("--->", err)
		return
	}
	fmt.Println(jsonstr)
	// 4、取0号对象，其余为js对象
	obj := objs[0]
	// 5、找关键词的路径

	action := &selm2.Action{Info: "新建服务"}
	if ok := action.Search(obj.Children); !ok {
		fmt.Println("unable find key info in pagedata")
		return
	}
	// 将标签数组反序列、初始化其他标签路径数组
	out := action.Reverse().Selector().ToString()

	fmt.Println(out.Path)
	fmt.Println(out.Link)
	fmt.Println(out.First)
	fmt.Println(out.Button)
	fmt.Println(out.Input)
	fmt.Println(out.Table)
	// 将路径数组拼凑成字符串
	// out := action.ToString()
	// fmt.Println(out.Link)
	// fmt.Println(out.First)
	// fmt.Println(out.Path)
	// arr = append(arr, start)
	// str := selm2.Reverse(&arr)
	// fmt.Println(str)
	// 6、点击蜜罐管理(buttom和a标签路径,谁不为空点谁),  这些操作需要分别实现
	//    输入操作:找text标签,
	//    提交操作:wd。submit即可
	//    启动、关闭、删除、还原操作
	// e, err := wd.FindElement(selenium.ByXPATH, str)
	// if err != nil {
	// 	fmt.Println("unable find hoyneypot btm", err)
	// 	return
	// }
	// return e.Click()
	return
}
