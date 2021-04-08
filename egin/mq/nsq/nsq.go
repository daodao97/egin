package nsq

import (
	"errors"
	"time"

	"github.com/nsqio/go-nsq"
)

type Producer struct {
	nsqAddress []string
	instance   *nsq.Producer
}

func (p *Producer) init() {
	if p.instance != nil && p.instance.Ping() == nil {
		return
	}

	// TODO 多 nsq 负载
	producer, err := nsq.NewProducer(p.nsqAddress[0], nsq.NewConfig())
	if err != nil {
		panic(err)
	}
	producer.SetLogger(nil, 0)
	p.instance = producer
}

func (p *Producer) Publish(topic string, message string) error {
	p.init()
	if p.instance == nil {
		return errors.New("producer is nil")
	}

	if message == "" { // 不能发布空串，否则会导致error
		return nil
	}
	err := p.instance.Publish(topic, []byte(message)) // 发布消息
	return err
}

func (p *Producer) Stop() {
	p.instance.Stop()
}

type Handle struct {
	handle func(message string)
}

func (hm *Handle) HandleMessage(msg *nsq.Message) error {
	hm.handle(string(msg.Body))
	msg.Finish()
	return nil
}

type Consumer struct {
	lookupAddress string
	nsqAddress    []string
}

func (nc *Consumer) Consumer(topic string, channel string, handler Handle) {
	address := nc.lookupAddress
	cfg := nsq.NewConfig()
	cfg.LookupdPollInterval = time.Second          // 设置重连时间
	c, err := nsq.NewConsumer(topic, channel, cfg) // 新建一个消费者
	if err != nil {
		panic(err)
	}
	c.SetLogger(nil, 0) // 屏蔽系统日志
	// c.AddHandler(&handler) // 添加消费者接口
	c.AddConcurrentHandlers(&handler, 200)

	// 建立NSQLookupd连接
	if err := c.ConnectToNSQLookupd(address); err != nil {
		panic(err)
	}

	// 建立多个nsqd连接
	if err := c.ConnectToNSQDs(nc.nsqAddress); err != nil {
		panic(err)
	}

	// 建立一个nsqd连接
	// if err := c.ConnectToNSQD("127.0.0.1:4150"); err != nil {
	//  panic(err)
	// }
}
