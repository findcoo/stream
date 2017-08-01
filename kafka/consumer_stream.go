package kafka

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
)

type (
	// MessageHandler consumer가 메세지를 받을때, 실행된다.
	MessageHandler func(*sarama.ConsumerMessage) *sarama.ConsumerMessage
	// ErrorHandler consumer가 broker로 부터 err 응답을 받을때, 실행
	ErrorHandler func(error)
	// NotificationHandler consumer가 broker로 부터 알림을 받을때, 실행
	NotificationHandler func(*cluster.Notification)
)

// ConsumerGroup sarama cluster를 이용한 consumer
type ConsumerGroup struct {
	cancel   chan struct{}
	consumer *cluster.Consumer
	Handler  ConsumHandler
}

// ConsumHandler consum 이벤트에 맞춰 실행될 핸들러
type ConsumHandler struct {
	AtSubscribe MessageHandler
	AtError     ErrorHandler
	AtNotified  NotificationHandler
}

// NewConsumerGroup consumer group 생성
func NewConsumerGroup(groupName string, addrs, topics []string, handler ConsumHandler) *ConsumerGroup {
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true

	c, err := cluster.NewConsumer(addrs, groupName, topics, config)
	if err != nil {
		_ = c.Close()
		log.Fatal(err)
	}

	cg := &ConsumerGroup{
		cancel:   make(chan struct{}, 1),
		consumer: c,
		Handler:  handler,
	}

	return cg
}

// Subscribe broker로 부터 메세지를 구독한다.
func (cg *ConsumerGroup) Subscribe(call MessageHandler) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill, syscall.SIGTERM)

	if call != nil {
		cg.Handler.AtSubscribe = call
	}

SubLoop:
	for {
		select {
		case msg, more := <-cg.consumer.Messages():
			cg.Handler.AtSubscribe(msg)
			if more {
				cg.consumer.MarkOffset(msg, "")
			}
		case err := <-cg.consumer.Errors():
			cg.Handler.AtError(err)
		case ntf := <-cg.consumer.Notifications():
			cg.Handler.AtNotified(ntf)
		case <-sig:
			break SubLoop
		case <-cg.cancel:
			break SubLoop
		}
	}
}

// Cancel 구독을 취소시킨다.
func (cg *ConsumerGroup) Cancel() {
	cg.cancel <- struct{}{}
}
