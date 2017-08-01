package kafka

import (
	"log"

	"github.com/Shopify/sarama"
	"github.com/findcoo/stream"
)

// ProducerStream sarama AsyncProducer를 이용한 stream 구조체
type ProducerStream struct {
	stream       chan *sarama.ProducerMessage
	producer     sarama.AsyncProducer
	Observer     *stream.Observer
	Subscribable func(*sarama.ProducerMessage)
}

// NewProducerStream stream 생성
func NewProducerStream(addr string) *ProducerStream {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll

	p, err := sarama.NewAsyncProducer([]string{addr}, config)
	if err != nil {
		log.Fatal(err)
	}

	ps := &ProducerStream{
		stream:   make(chan *sarama.ProducerMessage, 1),
		producer: p,
		Observer: stream.NewObserver(),
	}

	return ps
}

// Send ProducerMessage를 broker로 보낸다.
func (ps *ProducerStream) Send(data *sarama.ProducerMessage) {
	ps.stream <- data
}

// Publish stream 시작하고 stream을 반환한다.
func (ps *ProducerStream) Publish() *ProducerStream {
	ps.Observer.Observ()

	return ps
}

// Subscribe 데이터 구독
func (ps *ProducerStream) Subscribe(callArray ...func(*sarama.ProducerMessage)) {
SubLoop:
	for {
		select {
		case <-ps.Observer.DoneSubscribe:
			break SubLoop
		case data := <-ps.stream:
			if len(callArray) != 0 {
				for _, call := range callArray {
					call(data)
				}
			} else {
				ps.Subscribable(data)
				ps.producer.Input() <- data
			}
		}
	}
}

// Cancel stream 진행을 취소시킨다.
func (ps *ProducerStream) Cancel() {
	ps.Observer.Cancel()
}
