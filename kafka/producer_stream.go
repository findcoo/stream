package kafka

import (
	"log"

	"github.com/Shopify/sarama"
	"github.com/findcoo/stream"
)

// ProceedMessageHandler 메세지가 broker로 전달된 후 실행되는 핸들러
type ProceedMessageHandler func(*sarama.ProducerMessage)

// ProducerStream sarama AsyncProducer를 이용한 stream 구조체
type ProducerStream struct {
	stream    chan *sarama.ProducerMessage
	Producer  sarama.AsyncProducer
	AfterSend ProceedMessageHandler
	ErrFrom   stream.ErrHandler
	*stream.Observer
}

// NewProducerStream stream 생성
func NewProducerStream(addrs []string, obv *stream.Observer) *ProducerStream {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll

	p, err := sarama.NewAsyncProducer(addrs, config)
	if err != nil {
		log.Fatal(err)
	}

	ps := &ProducerStream{
		stream:    make(chan *sarama.ProducerMessage, 1),
		Producer:  p,
		Observer:  obv,
		AfterSend: func(*sarama.ProducerMessage) {},
		ErrFrom:   func(error) {},
	}

	return ps
}

// Send ProducerMessage를 Subscribable에 전달
func (ps *ProducerStream) Send(msg *sarama.ProducerMessage) {
	ps.WG.Add(1)
	ps.stream <- msg
}

// Publish Observable의 데이터를 구독한 후 broker로 메세지를 전송한다.
func (ps *ProducerStream) Publish(target func()) {
	ps.Watch(target)

SubLoop:
	for {
		select {
		case msg := <-ps.stream:
			ps.Producer.Input() <- msg
			ps.AfterSend(msg)
		case <-ps.Producer.Successes():
			ps.WG.Done()
		case err := <-ps.Producer.Errors():
			ps.ErrFrom(err)
		case <-ps.DoneSubscribe:
			break SubLoop
		}
	}
}
