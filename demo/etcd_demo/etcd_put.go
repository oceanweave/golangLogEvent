package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
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
	value := `[{"path":"/Users/dufengyang/代码/log-event/demo/xxx.log","topic":"web_log"},{"path":"/Users/dufengyang/代码/log-event/demo/xxx1.log","topic":"web_log"},{"path":"/Users/dufengyang/代码/log-event/demo/xxx2.log","topic":"web_log"}]`
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//_, err = cli.Put(ctx, "/logagent/192.168.8.199/collect_config", value) // collect_log_key=/logagent/%s/collect_config
	_, err = cli.Put(ctx, "/logagent/192.168.8.199/collect_config", value) // collect_log_key=/logagent/%s/collect_config
	cancel()
	if err != nil {
		fmt.Printf("put to etcd failed, err: %v\n", err)
		return
	}
}
