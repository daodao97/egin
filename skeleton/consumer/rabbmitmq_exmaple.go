package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/daodao97/egin/egin/mq/rabbitmq"
)

var mode = flag.String("mode", "", "运行模式")

func main() {
	flag.Parse()
	switch *mode {
	case "push":
		push()
	case "consume":
		consume()
	default:
		log.Fatal("place set mode")
	}
}

func push() {
	mq := rabbitmq.RabbitMQ{
		Connection: "default",
		QueueName:  "egin_test_queue",
	}
	err := mq.Publish("this is a test message")
	if err == nil {
		fmt.Println("投递成功")
	} else {
		log.Fatalf("投递失败, 原因: %s", err)
	}
}

func consume() {
	mq := rabbitmq.RabbitMQ{
		Connection: "default",
		QueueName:  "egin_test_queue",
		ConsumeConf: rabbitmq.ConsumeOptions{
			ConsumerName: "my_consume",
		},
	}
	mq.Consume(func(msg string) error {
		log.Printf("Received a message: %s", msg)
		return nil
	})
}
