package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"log_transfer/conf"
	"log_transfer/es"
	"log_transfer/kafka"
)

// log transfer
// 将日志数据从kafka取出来 发往es

func main() {
	// 0. 加载配置文件
	var cfg = new(conf.LogTransferCfg)         // 返回的是指针
	err := ini.MapTo(cfg, "./conf/config.ini") // 传指针  go语言都是值传递 就是复制 因此要传指针
	if err != nil {
		fmt.Printf("init config failed, err:%v\n", err)
		return
	}
	fmt.Printf("cfg:%v\n", cfg)
	// 1. 初始化 ES  传入连接地址 缓存通道大小 go程数目 将数据发往ES
	err = es.Init(cfg.ESConf.Address, cfg.ESConf.ChanSize, cfg.ESConf.NumsGo)
	if err != nil {
		fmt.Printf("init ES failed, err:%v\n", err)
		return
	}
	fmt.Println("init es success.")
	// 2. 初始化kafka
	// 2.1 连接kafka，穿件分区的消费者
	// 2.2 每个分区的消费者分别取出数据，通过es.SendToES() 将数据发往es
	err = kafka.Init([]string{cfg.KafkaConf.Address}, cfg.KafkaConf.Topic)
	if err != nil {
		fmt.Printf("init kafka consumer failed, err:%v\n", err)
		return
	}
	select {}
}
