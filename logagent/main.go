package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"logagent/conf"
	"logagent/kafka"
	"logagent/taillog"
	"time"
)

var (
	cfg = new(conf.AppConf) // 注意此处 要指针 因为以后要修改 所以要传指针
)

func run() {
	// 1. 读取日志
	for {
		select {
		case line := <-taillog.ReadcChan():
			// 2. 发送到 kafka
			kafka.SendToKafka(cfg.Topic, line.Text)
		default:
			time.Sleep(time.Second)
		}
	}

}

func main() {
	// 0. 加载配置文件
	// ini 文档 https://ini.unknwon.io/docs/advanced/map_and_reflect
	err := ini.MapTo(cfg, "./conf/config.ini") // 传指针
	// 1. 初始化 kafka 连接
	err = kafka.Init([]string{cfg.Address})
	if err != nil {
		fmt.Printf("Init Kafka failed, err: %v\n", err)
		return
	}
	fmt.Println("Init Kafka success.")
	// 2. 打开日志文件 准备收集日志
	err = taillog.Init(cfg.FileName)
	if err != nil {
		fmt.Printf("Init taillog failed, err: %v\n", err)
		return
	}
	fmt.Println("Init Taillog success.")
	run()
}
