package kafka

import (
	"fmt"

	"github.com/golang/glog"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"

	"github.com/daodao97/egin/utils/config"
)

// 目前仅支持最简单的内网模式
// 查看源码 example 解锁更多姿势

// kafka的消费者
func NewConsumer(kafkaConnection string, groupId string, subTopicS []string, handler func(message string), maxExecuteMessageCount int) {
	server, ok := config.Config.Kafka[kafkaConnection]
	if !ok {
		panic(fmt.Sprintf("kafka connection [%s] not found", kafkaConnection))
	}
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": server,
		"group.id":          groupId,
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}

	err = c.SubscribeTopics(subTopicS, nil)

	if err != nil {
		panic(err)
	}

	counter := 0

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			// fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			go handler(string(msg.Value))
		} else {
			// The client will automatically try to recover from all errors.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
		counter++
		if maxExecuteMessageCount > 0 && counter >= maxExecuteMessageCount {
			break
		}
	}

	_ = c.Close()
}

// kafka 生产者
func NewProducer(kafkaConnection string, messages []string, topic string) error {
	server, ok := config.Config.Kafka[kafkaConnection]
	if !ok {
		panic(fmt.Sprintf("kafka connection [%s] not found", kafkaConnection))
	}
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": server})
	if err != nil {
		panic(err)
	}

	defer p.Close()

	deliveryChan := make(chan kafka.Event)

	// Produce messages to topic (asynchronously)
	for _, message := range messages {
		err = p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(message),
		}, deliveryChan)

		if err != nil {
			glog.Info("kafka produce msg error")
		}
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if err = m.TopicPartition.Error; err != nil {
		fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
	} else {
		fmt.Printf("Delivered message to topic %s [%d] at offset %v\n", *m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}

	close(deliveryChan)

	return err
}
