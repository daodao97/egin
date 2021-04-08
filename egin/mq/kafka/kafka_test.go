package kafka

import (
	"testing"

	"github.com/go-playground/assert/v2"

	"github.com/daodao97/egin/utils/config"
)

func init() {
	config.Config.Kafka = make(map[string]string)
	config.Config.Kafka["default"] = "localhost:9092"
}

func TestNewProducer(t *testing.T) {
	messages := []string{"hi"}
	err := NewProducer("default", messages, "test")
	assert.Equal(t, nil, err)
}

func TestNewConsumer(t *testing.T) {
	messages := []string{"hi", "hi"}
	_ = NewProducer("default", messages, "test")
	handler := func(message string) {
		assert.Equal(t, "hi", message)
	}
	NewConsumer("default", "egin", []string{"test"}, handler, 1)
	t.Log("this is over")
}
