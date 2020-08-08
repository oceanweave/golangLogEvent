package conf

type LogTransferCfg struct {
	KafkaConf `ini:"kafka"` // 按节加载
	ESConf    `ini:"es"`
}

type KafkaConf struct {
	Address string `ini:"address"`
	Topic   string `ini:"topic"`
}

type ESConf struct {
	Address  string `ini:"address"`
	ChanSize int    `ini:"chansize"`
	NumsGo   int    `ini:"numsgo"`
}
