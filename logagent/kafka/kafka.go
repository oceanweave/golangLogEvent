package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
)

// 专门往 kafka 里写日志的模块

var (
	client sarama.SyncProducer // 声明一个全局的连接 kafka 的生产者 client
)

func Init(addrs []string) (err error) {
	config := sarama.NewConfig()
	// 发送数据需要 leader 和 follow 都确认
	config.Producer.RequiredAcks = sarama.WaitForAll
	// 新选出一个 partition Random 轮询
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	// 成功交付的消息将在 success channel 返回
	config.Producer.Return.Successes = true

	// 连接kafka
	client, err = sarama.NewSyncProducer(addrs, config)
	if err != nil {
		fmt.Println("producer closed, err:", err)
	}
	// defer client.Close()
	return err
}

func SendToKafka(topic, data string) {
	// 构造一个消息
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.StringEncoder(data)
	// 发送消息
	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		fmt.Println("send msg failed, err:", err)
		return
	}
	fmt.Printf("pid: %v offset:%v \n", pid, offset)
}
