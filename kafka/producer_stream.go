package kafka

import (
	"log"

	"github.com/Shopify/sarama"
	"github.com/findcoo/stream"
)

// ProducerStream sarama AsyncProducer를 이용한 stream 구조체
type ProducerStream struct {
	stream      chan *sarama.ProducerMessage
	producer    sarama.AsyncProducer
	Observer    *stream.Observer
	AtSubscribe func(*sarama.ProducerMessage)
}

// NewProducerStream stream 생성
func NewProducerStream(addr string, handler stream.ObservHandler) *ProducerStream {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll

	p, err := sarama.NewAsyncProducer([]string{addr}, config)
	if err != nil {
		log.Fatal(err)
	}

	ps := &ProducerStream{
		stream:   make(chan *sarama.ProducerMessage, 1),
		producer: p,
		Observer: stream.NewObserver(handler),
	}

	return ps
}

// Send ProducerMessage를 Subscribable에 전달
func (ps *ProducerStream) Send(data *sarama.ProducerMessage) {
	ps.stream <- data
}

// Publish Observable의 데이터를 구독한 후 broker로 메세지를 전송한다.
func (ps *ProducerStream) Publish(call func(*sarama.ProducerMessage)) {
	if call != nil {
		ps.AtSubscribe = call
	}

	ps.Observer.Observ()

	// NOTE Observable 구독 루프
SubLoop:
	for {
		select {
		case <-ps.Observer.DoneSubscribe:
			break SubLoop
		case data := <-ps.stream:
			ps.AtSubscribe(data)

			select {
			case ps.producer.Input() <- data:
			case err := <-ps.producer.Errors():
				log.Println("Failed to produce message", err)
			}
		}
	}
}

// Cancel stream 진행을 취소시킨다.
func (ps *ProducerStream) Cancel() {
	ps.Observer.Cancel()
}
