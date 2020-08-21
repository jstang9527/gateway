package selm

import (
	"fmt"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

const (
	seleniumPath = "/root/chromedriver"
	port         = 9515
)
const (
	entry = iota + 1
	input
	click
	monitor
	parse
	custom
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
	service *selenium.Service
	caps    = selenium.Capabilities{"browserName": "chrome"}
	ops     = []selenium.ServiceOption{}
)

// Init 创建一个selenium后台服务
func Init() (err error) {
	service, err = selenium.NewChromeDriverService(seleniumPath, port, ops...)
	if err != nil {
		return
	}
	fmt.Println("Start Selenium Service Listen on :::", port)
	return
}

// Stop 停止服务
func Stop() {
	fmt.Println("\n[INFO] Selenium ServerStop stopped")
	service.Stop()
}
