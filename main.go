package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jstang9527/gateway/thirdpart/es"

	"github.com/jstang9527/gateway/thirdpart/selm"

	"github.com/jstang9527/gateway/router"

	"github.com/e421083458/golang_common/lib"
)

func main() {
	lib.InitModule("./conf/dev/", []string{"base", "mysql", "redis"})
	defer lib.Destroy()
	// 启动selenium后台服务
	if err := selm.Init(); err != nil {
		fmt.Println("Failed Init Selenium Service, Info: ", err)
		return
	}
	// 创建ES连接池、启动推送守护进程
	if err := es.InitES(); err != nil {
		fmt.Println("Failed Init Elasticsearch Service, Info: ", err)
		return
	}

	router.HTTPServerRun()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	selm.Stop()
	es.Stop()
	router.HTTPServerStop()
}
