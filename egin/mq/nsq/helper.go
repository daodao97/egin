package nsq

import (
	"github.com/daodao97/egin/egin/utils/config"
)

// 创建nsq消费者
func NewConsumer(connection string) *Consumer {
	conf, ok := config.Config.Nsq[connection]
	if !ok {
		panic("nsq conf not found")
	}
	return &Consumer{nsqAddress: conf.NsqAddress, lookupAddress: conf.LookupAddress}
}

// 创建nsq生产者
func NewProducer(connection string) *Producer {
	conf, ok := config.Config.Nsq[connection]
	if !ok {
		panic("nsq conf not found")
	}
	return &Producer{nsqAddress: conf.NsqAddress}
}
