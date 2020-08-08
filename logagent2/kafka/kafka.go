package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
)

// 专门往 kafka 里写日志的模块

var (
	client      sarama.SyncProducer // 声明一个全局的连接 kafka 的生产者 client
	logDataChan chan *logData       // 内部有缓冲的channel 用于异步通信 外部数据传送过来
)

type logData struct {
	topic string
	data  string
}

func Init(addrs []string, maxSize int) (err error) {
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
	// 初始化 logDataChan 通道
	logDataChan = make(chan *logData, maxSize)
	// 开始后台的goroutine从通道中取出数据发往kafka
	go sendToKafka()
	return
}

// 给外部暴露的一个函数， 该函数至把日志数据发动到一个内部的channel中
func SendToChan(topic, data string) {
	msg := &logData{
		topic: topic,
		data:  data,
	}
	logDataChan <- msg
}

// 真正网 kafka 发送日志的函数
func sendToKafka() {
	for {
		select {
		case ld := <-logDataChan: // 逐条取出 logData 消息
			// 构造一个消息
			msg := &sarama.ProducerMessage{}
			msg.Topic = ld.topic
			msg.Value = sarama.StringEncoder(ld.data)
			// 发送消息
			pid, offset, err := client.SendMessage(msg)
			if err != nil {
				fmt.Println("send msg failed, err:", err)
				return
			}
			fmt.Printf("pid: %v offset:%v \n", pid, offset)
		}
	}
}
