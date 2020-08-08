package taillog

import (
	"fmt"
	"logagent/etcd"
	"time"
)

var tskMgr *tailLogMgr

// tailTask 管理者
type tailLogMgr struct {
	logEntrys   []*etcd.LogEntry
	tskMap      map[string]*TailTask
	newConfChan chan []*etcd.LogEntry
}

func Init(logEntryConf []*etcd.LogEntry) {
	tskMgr = &tailLogMgr{
		logEntrys:   logEntryConf, // 把当前的日志收集箱配置信息保存起来
		tskMap:      make(map[string]*TailTask, 16),
		newConfChan: make(chan []*etcd.LogEntry), // 无缓冲的通道 通道初始化
	}
	// 3.1 循环每一个日志收集项， 创建TailObj
	for _, logEntry := range logEntryConf {
		// 初始化是 创建了多少分 tailtask(就是多少个日志文件） 都要记录下来  为了方便后序判断该任务状态 增删改
		tailObj := NewTailTask(logEntry.Path, logEntry.Topic)
		mk := fmt.Sprintf("%s_%s", logEntry.Path, logEntry.Topic)
		tskMgr.tskMap[mk] = tailObj
	}

	go tskMgr.run()
}

// 监听自己的 newConfChan 有了新的配置过来之后就做对应的处理

func (t *tailLogMgr) run() {
	for {
		select {
		case newConf := <-t.newConfChan:
			for _, conf := range newConf {
				mk := fmt.Sprintf("%s_%s", conf.Path, conf.Topic)
				_, ok := t.tskMap[mk]
				if ok {
					// 原来就有， 不需要操作
					continue
				} else {
					// 新增的
					tailObj := NewTailTask(conf.Path, conf.Topic)
					t.tskMap[mk] = tailObj
				}
			}
			// 找出原来 t.logEntrys 有，但现在newConf中没有的，要删掉
			for _, c1 := range t.logEntrys {
				isDelete := true
				for _, c2 := range newConf {
					if c1.Path == c2.Path && c1.Topic == c2.Topic {
						isDelete = false
						continue
					}
				}
				if isDelete {
					// 把对应的tailObj给停掉
					mk := fmt.Sprintf("%s_%s", c1.Path, c1.Topic)
					t.tskMap[mk].cancelFunc() // 执行 cancel() ctx.Done()收到信号
				}
			}
			// 1. 配置新增
			// 2，配置删除
			// 3. 配置变更
			//fmt.Println("新的配置来啦", newConf)
		default:
			time.Sleep(time.Second)
		}
	}
}

// 向外暴露 tskMgr 的newConfChan
func NewConfChan() chan<- []*etcd.LogEntry {
	return tskMgr.newConfChan
}
