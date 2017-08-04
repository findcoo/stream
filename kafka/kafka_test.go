package kafka

import (
	"log"
	"strconv"
	"testing"

	"github.com/Shopify/sarama"
)

func rangeProduce(m, n int, topic string) {
	h := &ProduceHandler{
		AfterSend: func(msg *sarama.ProducerMessage) {
			log.Print(msg)
		},
		ErrFrom: func(err error) {
			log.Print(err)
		},
	}

	p := NewProducerStream([]string{"localhost:9092"}, h, nil)
	p.Observer.Handler.AtComplete = func() {
		log.Print("end produce")
	}

	p.Observer.SetObservable(func() {
		for i := m; i <= n; i++ {
			msg := &sarama.ProducerMessage{
				Topic: topic,
				Key:   sarama.StringEncoder(strconv.Itoa(i)),
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

func rangeConsume(n int, topic string) {
	h := DefaultConsumHandler()
	h.AtError = func(err error) {
		log.Print(err)
	}

	c := NewConsumerGroup("test_group", []string{"127.0.0.1:9092"}, []string{topic}, h)

	var count int
	c.Subscribe(func(msg *sarama.ConsumerMessage) {
		log.Printf("offset: %d\n", msg.Offset)
		count++
		if count == n {
			c.Cancel()
		}
	})
}

func TestKafka(t *testing.T) {
	rangeProduce(0, 10, "test")
	rangeConsume(10, "test")
}
