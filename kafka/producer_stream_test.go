package kafka

import (
	"strconv"
	"testing"

	"github.com/Shopify/sarama"
)

func TestProducer(t *testing.T) {
	h := &ProduceHandler{
		AfterSend: func(msg *sarama.ProducerMessage) {
			t.Log(msg)
		},
		ErrFrom: func(error) {},
	}

	p := NewProducerStream("127.0.0.1:9092", h, nil)
	p.Observer.SetObservable(func() {
		for i := 0; i <= 10; i++ {
			msg := &sarama.ProducerMessage{
				Topic: "test",
				Key:   nil,
				Value: sarama.StringEncoder(strconv.Itoa(i)),
			}

			select {
			case <-p.Observer.AfterCancel():
				return
			default:
				p.Send(msg)
			}
		}
		p.Observer.OnComplete()
	})

	p.Publish()
}
