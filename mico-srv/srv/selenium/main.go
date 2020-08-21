package main

import (
	"fmt"
	"time"

	"github.com/jstang9527/gateway/mico-srv/modules"
	pb "github.com/jstang9527/gateway/mico-srv/srv/pb"
	mops "github.com/jstang9527/opentest/ops"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"golang.org/x/net/context"
)

const (
	address      = "127.0.0.1:50052"
	seleniumPath = "/root/chromedriver"
	port         = 9515
)

var (
	chromeCaps = chrome.Capabilities{
		Prefs: map[string]interface{}{"profile.managed_default_content_settings.images": 2},
		Path:  "",
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

// 定义seleniumService并实现约定的接口
type seleniumService struct{}

// RunTest 实现服务接口
func (ss seleniumService) RunTest(ctx context.Context, in *pb.SeleniumRequest) (*pb.SeleniumResponse, error) {
	resp := new(pb.SeleniumResponse)
	// 1.注册selenium客户端,生成webdriver控制器
	caps := selenium.Capabilities{"browserName": "chrome"}
	caps.AddChrome(chromeCaps)
	wd, err := selenium.NewRemote(caps, "http://127.0.0.1:9515/wd/hub")
	if err != nil {
		resp.Message = fmt.Sprintf("create webdriver failed, info: %v", err)
		return resp, nil
	}
	defer wd.Quit()
	wd.SetImplicitWaitTimeout(time.Second * time.Duration(in.SearchTimeout))
	// 2.执行测试
	tt := &modules.TestTask{URL: in.Url}
	go tt.Run(wd) //这是web controller调用的
	// tt.Run(wd) //这是模块测试调用的
	time.Sleep(time.Second * 2)
	return resp, nil
}

func main() {
	//1.开启selenium服务
	ops := []selenium.ServiceOption{}
	service, err := selenium.NewChromeDriverService(seleniumPath, port, ops...)
	if err != nil {
		fmt.Printf("Error starting the ChromeDriver server: %v", err)
	}
	defer service.Stop()
	mops.Console("Selinium Service Listen on :::%v", port)
	// 2.启GRPC 服务
	// listen, err := net.Listen("tcp", address)
	// if err != nil {
	// 	fmt.Printf("Failed to listen: %v\n", err)
	// 	return
	// }

	// s := grpc.NewServer()
	// pb.RegisterSeleniumServer(s, seleniumService{})
	// mops.Console("Selenium gRPC Service Listen on :::%v", address)
	// s.Serve(listen)
}
