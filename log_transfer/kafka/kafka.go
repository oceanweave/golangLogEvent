package kafka

import (
	//"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"log_transfer/es"
)

// 初始化kafka消费者 从kafka取数据发往ES

// LogData...
//type LogData struct {//	Topic string `json:"topic"`
//	Data  string `json:"data"`
//}

// Init初始化
func Init(addrs []string, topic string) (err error) {
	// 连接kafka
	consumer, err := sarama.NewConsumer(addrs, nil)
	if err != nil {
		fmt.Printf("fail to start sonsumer, err:%v\n", err)
		return err
	}
	// 根据topic取到所有的分区
	partitionList, err := consumer.Partitions(topic)
	if err != nil {
		fmt.Printf("fail to get list of partition, err: %v\n", err)
		return err
	}
	fmt.Println("connet to kafka success.")
	fmt.Println("分区列表：", partitionList)
	for partition := range partitionList { // 遍历所有分区
		// 针对每个分区创建一个对应的分区消费者
		pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("failed to start consumer for partition %d, err: %v \n", partition, err)
			return err
		}
		//defer pc.AsyncClose()
		fmt.Println("进入for 分区遍历")
		// 异步从每个分区消费信息
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				fmt.Printf("Partition:%d Offset:%d Key:%v Value:%v \n", msg.Partition, msg.Offset, msg.Key, string(msg.Value))
				// 直接发往ES
				//var ld = new(LogData)
				//ld := map[string]interface{}{
				//	"data": msg.Value,
				//}
				ld := es.LogData{Topic: topic, Data: string(msg.Value)}
				es.SendToESChan(&ld)
				//fmt.Println(ld)

				//err := json.Unmarshal(msg.Value, ld) // 将json 数据解析到 ld 中
				//if err != nil {
				//	fmt.Printf("unmarshal failed, err:%v\n", err)
				//	continue
				//}
				//tmp, _ := json.Marshal(ld)
				//es.SendToES(topic, ld)
			}
		}(pc)
	}
	//select {} // 这个地方需要阻塞住 要不然退出 会导致上面 go程退出 无法读取kafka分区中数据
	return err
}
