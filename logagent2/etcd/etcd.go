package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

var (
	cli *clientv3.Client
)

// 需要收集的日志的配置信息
type LogEntry struct {
	Path  string `json:"path"`  // 日志存放的路径
	Topic string `json:"topic"` // 日志要发往kafka中的拿个topic
}

// 初始化 ETCD 的函数
func Init(addr string, timeout time.Duration) (err error) {
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{addr},
		DialTimeout: timeout,
	})
	if err != nil {
		fmt.Printf("connect to etcd failed, err: %v\n", err)
		return
	}
	fmt.Println("connect to etcd success")
	return
}

// 从 ETCD中能根据key 获取配置项
func GetConf(key string) (logEntryConf []*LogEntry, err error) {
	// get
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := cli.Get(ctx, key)
	cancel()
	if err != nil {
		fmt.Printf("get from etcd failed, err: %v\n", err)
		return
	}
	//fmt.Println("开始解析")
	for _, ev := range resp.Kvs {
		// 获取 etcd 中key 对应的全部 json 数据
		//fmt.Printf("%s:%s\n", ev.Key, ev.Value)
		err = json.Unmarshal(ev.Value, &logEntryConf) // 解析到 logEntryConf中
		if err != nil {
			fmt.Printf("unmarshal etcd value failed, err:%v\n", err)
		}
		// 遍历json数据 获取path-topic对
		for index, value := range logEntryConf {
			fmt.Printf("index: %v, value: %v\n", index, value)
		}
	}
	fmt.Println("get from etcd success.")
	return
}

// etcd watch
func WatchConf(key string, newConfCh chan<- []*LogEntry) {
	// watch key: dfy change  watch 一直在等
	rch := cli.Watch(context.Background(), key) //  <-chan WatchResponse 返回通道
	// 从通道尝试取值 监视的信息
	for wresp := range rch {
		for _, evt := range wresp.Events {
			fmt.Printf("Type: %s Key: %s Value: %s\n", evt.Type, evt.Kv.Key, evt.Kv.Value)
			// 通知 taillog.tskMgr
			// 1. 先判断操作的类型
			var newConf []*LogEntry
			if evt.Type != clientv3.EventTypeDelete {
				// 不是删除操作的话解析出来
				err := json.Unmarshal(evt.Kv.Value, &newConf)
				if err != nil {
					fmt.Printf("json unmarshal failed, err: %v", err)
					continue
				}
			}
			// 如果是删除操作 手动传递一个空的配置项
			fmt.Printf("get new conf: %v\n", newConf)
			newConfCh <- newConf
		}
	}
}
