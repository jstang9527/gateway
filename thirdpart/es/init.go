package es

import (
	"context"
	"fmt"
	"time"

	"github.com/e421083458/golang_common/lib"
	"github.com/jstang9527/gateway/dto"
	"github.com/olivere/elastic/v7"
)

var client *elastic.Client
var index string //索引,相当于数据库
var ch chan *dto.TestItem

// InitES 初始化
func InitES() (err error) {
	esip := lib.GetStringConf("base.es.host")
	esport := lib.GetIntConf("base.es.port")
	esindex := lib.GetStringConf("base.es.index")
	chansize := lib.GetIntConf("base.es.chansize")
	nums := lib.GetIntConf("base.es.nums")
	url := fmt.Sprintf("http://%s:%d", esip, esport)
	client, err = elastic.NewClient(elastic.SetURL(url))
	if err != nil {
		return
	}
	index = esindex
	// 通道初始化
	ch = make(chan *dto.TestItem, chansize)
	// PushES异步协程数量
	for i := 0; i < nums; i++ {
		go sendESDaemon()
	}
	fmt.Printf("[info] Start ElasticSearch Service Listen On :::%v\n", esport)
	fmt.Printf("[info] Build %v Channal, Max Containers: %v\n", nums, chansize)
	return
}

// Stop ...
func Stop() {
	defer func() {
		fmt.Println("Stopping Elasticsearch Service...")
	}()
	client.Stop()
}

//从管道拿数据发给ES, 这个初始化就要进行了(后台进程)
func sendESDaemon() {
	// put1, err := client.Index().Index(indexStr).Type(typeStr).BodyJson(data).Do(context.Background())
	fmt.Println("Start ES Channal ...")
	for {
		select {
		case item := <-ch:
			put1, err := client.Index().Index(index).BodyJson(item).Do(context.Background()) //BodyJson必须为结构体
			if err != nil {
				fmt.Printf("failed push data to ES from channal, err: %v\n", err)
			}
			fmt.Printf("indexed %s %s to index %s, type %s, data->%v\n", item.ProjectID, put1.Id, put1.Index, put1.Type, item.Message)
		default:
			time.Sleep(time.Second)
		}
	}
}
