package rabbitmq

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"

	"github.com/daodao97/egin/utils/config"
	"github.com/daodao97/egin/utils/logger"
)

var mqLogger = logger.NewLogger("mq")

type QueueOptions struct {
	QueueName  string     // queue 名称
	Durable    bool       // 是否持久化
	AutoDelete bool       // 是否为自动删除
	Exclusive  bool       // 是否具有排他性
	NoWait     bool       // 是否阻塞
	Args       amqp.Table // 额外属性
}

// queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args Table
type ConsumeOptions struct {
	ConsumerName string     // 用于区分不同的消费者
	AutoAck      bool       // 是否自动应答
	Exclusive    bool       // 是否具有排他性
	NoLocal      bool       // 如果设置为true，表示不能将同一个connection中发送的消息传递给这个connection中的消费者
	NoWait       bool       // 队列消费是否阻塞
	Args         amqp.Table // 额外属性
}

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	// 队列名称
	QueueName string
	// 交换机
	Exchange string
	// Key
	Key string
	// 配置中的链接名
	Connection  string
	QueueConf   QueueOptions
	ConsumeConf ConsumeOptions
}

func (r *RabbitMQ) init() {
	if r.conn != nil {
		return
	}
	if r.Connection == "" {
		r.Connection = "default"
	}

	conf := config.Config.RabbitMQ[r.Connection]

	var mqUrl = fmt.Sprintf("amqp://%s:%s@%s:%d/%s", conf.User, conf.Passwd, conf.Host, conf.Port, conf.Vhost)

	conn, err := amqp.Dial(mqUrl)
	if err != nil {
		mqLogger.Error(fmt.Sprintf("rabbitmq dial fail [%s:%s]", r.Connection, err))
		return
	}
	r.conn = conn

	channel, err := r.conn.Channel()
	if err != nil {
		mqLogger.Error(fmt.Sprintf("rabbitmq channel fail [%s:%s]", r.Connection, err))
		return
	}
	r.channel = channel
}

func (r *RabbitMQ) Queue() *RabbitMQ {
	r.init()
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		r.QueueConf.Durable,
		r.QueueConf.AutoDelete,
		r.QueueConf.Exclusive,
		r.QueueConf.NoWait,
		r.QueueConf.Args,
	)
	if err != nil {
		mqLogger.Error(fmt.Sprintf("rabbitmq QueueDeclare fail [%s:%s]", r.Connection, err))
	}

	return r
}

func (r *RabbitMQ) Publish(message string) error {
	r.Queue()
	return r.channel.Publish(
		r.Exchange,
		r.QueueName,
		// 如果为true, 会根据exchange类型和routkey规则，如果无法找到符合条件的队列那么会把发送的消息返回给发送者
		false,
		// 如果为true, 当exchange发送消息到队列后发现队列上没有绑定消费者，则会把消息发还给发送者
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
}

// Consume 使用 goroutine 消费消息
func (r *RabbitMQ) Consume(handler func(msg string) error) {
	r.Queue()
	// 接收消息
	messageChan, err := r.channel.Consume(
		r.QueueName,
		r.ConsumeConf.ConsumerName,
		r.ConsumeConf.AutoAck, // FIXME
		r.ConsumeConf.Exclusive,
		r.ConsumeConf.NoLocal,
		r.ConsumeConf.NoWait,
		r.ConsumeConf.Args,
	)

	if err != nil {
		fmt.Println(err)
	}

	forever := make(chan bool)
	// 启用协程处理消息
	go func() {
		for d := range messageChan {
			err := handler(string(d.Body))
			if err != nil {
				log.Fatalf("message handle fail: %s", err)
			} else {
				// FIXME 批处理
				err := d.Ack(false)
				if err != nil {
					mqLogger.Error(fmt.Sprintf("rabbitmq consume fail [%s:%s]", r.Connection, err))
				}
			}
		}
	}()
	log.Printf("[*] Waiting for message, To exit press CTRL+C")
	<-forever
}

// Destroy 断开channel和connection
func (r *RabbitMQ) Destroy() {
	_ = r.channel.Close()
	_ = r.conn.Close()
}
