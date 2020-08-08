package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Printf("connext to etcd failed, err: %v\n", err)
		return
	}
	fmt.Println("connect to etcd success")
	defer cli.Close()
	// watch key: dfy change  watch 一直在等
	rch := cli.Watch(context.Background(), "dfy") //  <-chan WatchResponse 返回通道
	// 从通道尝试取值 监视的信息
	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Printf("Type: %s Key: %s Value: %s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}
}
