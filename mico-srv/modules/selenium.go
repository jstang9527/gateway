package modules

import (
	"fmt"
	"time"

	"github.com/tebeka/selenium"
)

// TestTask ...
type TestTask struct {
	URL string
}

// Run 从DB中一条条读，或者使用管道; 错误信息不要打印，全部存es
func (tt *TestTask) Run(wd selenium.WebDriver) {
	var err error
	if err = wd.Get("https://172.31.50.39:65000/"); err != nil {
		fmt.Println(err)
		return
	}
	we, err := wd.FindElement(selenium.ByXPATH, `//*[@id="root"]/div/div[2]/ul/li[3]/div[1]`)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err = we.Click(); err != nil {
		fmt.Println(err)
		return
	}
	we, err = wd.FindElement(selenium.ByXPATH, `//*[@id="/honeypot-manage$Menu"]/li[1]/div[5]`)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = we.Click()
	time.Sleep(20 * time.Second)

}
