# kafka
improve kafka's reactivity

* Producer
  observer watch the Observable and directly publish to broker
  there's no need to subscribe the stream

  ```go
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
		for i := 0; i <= 10; i++ {
			msg := &sarama.ProducerMessage{
				Topic: "test",
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
  ```
* Consumer
  consumer not have observer thus acting like ovserver to watch broker

  ```go
	h := DefaultConsumHandler()
	h.AtError = func(err error) {
		log.Print(err)
	}

	c := NewConsumerGroup("test_group", []string{"localhost:9092"}, []string{"test"}, h)

	var count int
	c.Subscribe(func(msg *sarama.ConsumerMessage) {
		log.Printf("offset: %d\n", msg.Offset)
		count++
		if count == 10 {
			c.Cancel()
		}
	})
  ```
