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
	stream   chan *sarama.ProducerMessage
	Producer sarama.AsyncProducer
	Observer *stream.Observer
	Handler  ProduceHandler
}

// ProduceHandler producer 활동중 실행되는 핸들러들
type ProduceHandler struct {
	AfterSend ProceedMessageHandler
	ErrFrom   stream.ErrHandler
}

// DefaultProduceHandler 기본 핸들러 설정
func DefaultProduceHandler() *ProduceHandler {
	return &ProduceHandler{
		AfterSend: func(*sarama.ProducerMessage) {},
		ErrFrom:   func(error) {},
	}
}

// NewProducerStream stream 생성
func NewProducerStream(addrs []string, handler *ProduceHandler, obvHandler *stream.ObservHandler) *ProducerStream {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll

	if handler == nil {
		handler = DefaultProduceHandler()
	}

	p, err := sarama.NewAsyncProducer(addrs, config)
	if err != nil {
		log.Fatal(err)
	}

	ps := &ProducerStream{
		stream:   make(chan *sarama.ProducerMessage, 1),
		Producer: p,
		Observer: stream.NewObserver(obvHandler),
		Handler:  *handler,
	}

	return ps
}

// Send ProducerMessage를 Subscribable에 전달
func (ps *ProducerStream) Send(msg *sarama.ProducerMessage) {
	ps.Observer.WG.Add(1)
	ps.stream <- msg
}

// Publish Observable의 데이터를 구독한 후 broker로 메세지를 전송한다.
func (ps *ProducerStream) Publish() {
	ps.Observer.Observ()

SubLoop:
	for {
		select {
		case msg := <-ps.stream:
			ps.Producer.Input() <- msg
			ps.Handler.AfterSend(msg)
		case <-ps.Producer.Successes():
			ps.Observer.WG.Done()
		case err := <-ps.Producer.Errors():
			ps.Handler.ErrFrom(err)
		case <-ps.Observer.DoneSubscribe:
			break SubLoop
		}
	}
}

// Cancel stream 진행을 취소시킨다.
func (ps *ProducerStream) Cancel() {
	ps.Observer.Cancel()
}
