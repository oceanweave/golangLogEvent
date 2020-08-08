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
		fmt.Printf("connect to etcd failed, err: %v\n", err)
		return
	}
	fmt.Println("connect to etcd success")
	defer cli.Close()
	// put
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = cli.Put(ctx, "dfy", "test1")
	cancel()
	if err != nil {
		fmt.Printf("put to etcd failed, err: %v\n", err)
		return
	}
	fmt.Println("put data to etcd success")
	// get
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, "xxx")
	cancel()
	if err != nil {
		fmt.Printf("get from etcd failed, err: %v\n", err)
		return
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s:%s\n", ev.Key, ev.Value)
	}
	fmt.Println("get from etcd success.")
}
