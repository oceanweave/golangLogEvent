package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"logagent/conf"
	"logagent/etcd"
	"logagent/kafka"
	"logagent/taillog"
	"logagent/utils"
	"sync"
	"time"
)

var (
	cfg = new(conf.AppConf) // 注意此处 要指针 因为以后要修改 所以要传指针
)

//func run() {
//	// 1. 读取日志
//	for {
//		select {
//		case line := <-taillog.ReadcChan():
//			// 2. 发送到 kafka
//			kafka.SendToKafka(cfg.Topic, line.Text)
//		default:
//			time.Sleep(time.Second)
//		}
//	}
//
//}

func main() {
	// 0. 加载配置文件
	// ini 文档 https://ini.unknwon.io/docs/advanced/map_and_reflect
	err := ini.MapTo(cfg, "./conf/config.ini") // 传指针 cfg 为指针 这样才能修改
	// 验证配置是否正确读取
	//fmt.Println(cfg.EtcdConf.Address)
	//fmt.Println(cfg.KafkaConf.Address)
	//fmt.Println(cfg.EtcdConf.Timeout)

	// 1. 初始化 kafka 连接
	err = kafka.Init([]string{cfg.KafkaConf.Address}, cfg.KafkaConf.MaxSize)
	if err != nil {
		fmt.Printf("Init Kafka failed, err: %v\n", err)
		return
	}
	fmt.Println("Init Kafka success.")
	// 2. 初始化ETCD
	err = etcd.Init(cfg.EtcdConf.Address, time.Duration(cfg.EtcdConf.Timeout)*time.Second)
	if err != nil {
		fmt.Printf("Init ETCD failed, err: %v\n", err)
		return
	}
	fmt.Println("Init ETCD success.")
	// 2.1 从etcd 中获取日志收集项的配置信息
	// 为了实现每个logagent 都拉去自己独有的配置，所以要以自己的IP地址作为区分
	ipStr, err := utils.GetOutbounfIP()
	if err != nil {
		panic(err)
	}
	cfg.EtcdConf.Key = fmt.Sprintf(cfg.EtcdConf.Key, ipStr) //  第一个参数 包含了 %s /logagent/%s/collect_config
	logEntryConf, err := etcd.GetConf(cfg.EtcdConf.Key)
	if err != nil {
		fmt.Printf("etcd.GetConf failed, err:%v\n", err)
		return
	}
	fmt.Printf("get conf from etcd success, %v\n", logEntryConf)
	// 2.2 拍一个哨兵去监视日志收集箱的变化（有变化即使通知 logAgent 实现热加载配置）

	// 3. 收集日志发往kafka
	taillog.Init(logEntryConf)

	// taillog.NewConfChan() 中的返回的通道 需要现在Init 中初始化 因此放在后面 newConfChan
	newConfChan := taillog.NewConfChan() // 从 taillog 包中获取对外暴露的通道
	var wg sync.WaitGroup
	wg.Add(1)
	go etcd.WatchConf(cfg.EtcdConf.Key, newConfChan) // 哨兵发现最新的配置信息会通知上面的那个通道
	wg.Wait()

	//// 2. 打开日志文件 准备收集日志
	//err = taillog.Init(cfg.FileName)
	//if err != nil {
	//	fmt.Printf("Init taillog failed, err: %v\n", err)
	//	return
	//}
	//fmt.Println("Init Taillog success.")
	//run()
}
