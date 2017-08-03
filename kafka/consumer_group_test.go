package kafka

import (
	"log"
	"testing"

	"github.com/Shopify/sarama"
)

func TestConsumer(t *testing.T) {
	h := DefaultConsumHandler()
	h.AtError = func(err error) {
		t.Error(err)
	}

	c := NewConsumerGroup("test_group", []string{"127.0.0.1:9092"}, []string{"test"}, h)

	c.Subscribe(func(msg *sarama.ConsumerMessage) {
		log.Print(msg)
	})
}
