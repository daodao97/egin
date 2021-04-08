package nsq

import (
	"fmt"
	"testing"

	"github.com/go-playground/assert/v2"

	"github.com/daodao97/egin/egin/utils/config"
)

func init() {
	config.Config.Nsq = make(map[string]struct {
		LookupAddress string
		NsqAddress    []string
	})
	config.Config.Nsq["default"] = struct {
		LookupAddress string
		NsqAddress    []string
	}{
		LookupAddress: "127.0.0.1:4161",
		NsqAddress:    []string{"127.0.0.1:4152"},
	}
}

func TestNewProducer(t *testing.T) {
	producer := NewProducer("default")
	err := producer.Publish("test", "hi")
	assert.Equal(t, nil, err)
}

func TestNewConsumer(t *testing.T) {
	producer := NewProducer("default")
	_ = producer.Publish("test", "hi")

	handler := Handle{
		handle: func(message string) {
			fmt.Println(message)
		},
	}

	c := NewConsumer("default")

	c.Consumer("test", "test", handler)
}
