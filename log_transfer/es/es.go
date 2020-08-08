package es

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"strings"
	"time"
)

type LogData struct {
	Topic string `json:"topic"`
	Data  string `json:"data"`
}

var (
	client *elastic.Client
	ch     chan *LogData
)

// 初始化ES 准备接受kafka发过来的消息

func Init(addr string, chanSize, numsGo int) (err error) {
	if !strings.HasPrefix(addr, "http://") {
		addr = "http://" + addr
	}
	client, err = elastic.NewClient(elastic.SetURL(addr))
	if err != nil {
		// Handle error
		return err
	}
	// 先初始化通道
	ch = make(chan *LogData, chanSize)
	for i := 0; i < numsGo; i++ {
		go SendToES()
	}
	fmt.Println("connect to es success")
	return
}

// 外部发往 ES包 内部有缓冲channel 用于异步通信
func SendToESChan(msg *LogData) {
	ch <- msg
}

// SendToES 发送数据到ES  真正发送数据到ES
func SendToES() (err error) {
	for {
		select {
		case msg := <-ch:
			fmt.Println("send msg to es success")
			//ld := LogData{Data: data.(string)}
			// 链式操作 连接到指定index
			put1, err := client.Index().
				Index(msg.Topic).
				BodyJson(msg). // msg为可json化的结构体 注意此处要传入结构体
				Do(context.Background())
			if err != nil {
				// Handle error
				fmt.Println(err)
				continue
			}
			fmt.Printf("Indexed user %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
		default:
			time.Sleep(time.Second)
		}
	}

}
