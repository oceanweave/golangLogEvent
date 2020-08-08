package taillog

import (
	"context"
	"fmt"
	"github.com/hpcloud/tail"
	"logagent/kafka"
)

var (
	tailObj *tail.Tail // 全局
)

// TailTask 一个日志收集的任务
type TailTask struct {
	path     string
	topic    string
	instance *tail.Tail
	// 为了能实现退出 t.run()
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewTailTask(path, topic string) (tailObj *TailTask) {
	ctx, cancel := context.WithCancel(context.Background())
	tailObj = &TailTask{
		path:       path,
		topic:      topic,
		ctx:        ctx,
		cancelFunc: cancel,
	}
	tailObj.init() // 根据路径去打开对应的日志 给 t.instance 赋值
	return
}

func (t *TailTask) init() {
	// 每一个路径 创建一个对应到的 tailObj
	config := tail.Config{
		ReOpen:    true,                                 // 重新打开
		Follow:    true,                                 // 是否跟随
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 从文件的哪个地方开始读
		MustExist: false,                                // 文件不存在不报错
		Poll:      true,
	}
	var err error
	t.instance, err = tail.TailFile(t.path, config) // 初始化
	if err != nil {
		fmt.Println("tail file failed, err:", err)
		return
	}
	// 当goroutine执行的函数退出时，goroutine就结束了
	go t.run() // 写入kafka
	return
}

// 一直在收集日志
func (t *TailTask) run() {
	for {
		select {
		case <-t.ctx.Done(): // 退出
			fmt.Printf("tail task:%s_%s 结束了...\n", t.path, t.topic)
			return
		case line := <-t.instance.Lines: // 从tailObj的通道中一行一行读取日志数据
			// 3.2 发往kafka
			// kafka.SendToKafka(t.topic, line.Text) // 函数调用函数
			// 同步变为异步
			// 先把日志数据发到一个通道汇总
			kafka.SendToChan(t.topic, line.Text)
			// kafka  那个包中当杜的goroutine 去取日志数据发送到 kafka
		}
	}
}
