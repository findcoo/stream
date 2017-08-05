# stream
[![Go Travis Ci](https://travis-ci.org/findcoo/stream.svg?branch=master)](https://travis-ci.org/findcoo/stream)
[![Go Report Card](https://goreportcard.com/badge/github.com/findcoo/stream)](https://goreportcard.com/report/github.com/findcoo/stream)

make simple ReactiveX pattern for golang

# Subsets
* [Kafka](https://github.com/findcoo/stream/tree/master/kafka#kafka)

# Example
* bytes stream
  ```go
	handler := &ObservHandler{
		AtCancel:   func() {},
		AtComplete: func() {},
		AtError:    func(error) {},
	}

	stream := NewBytesStream(handler)
	stream.Observer.SetObservable(func() {
		for i := 0; i <= 10; i++ {
			stream.Send([]byte{byte(i)})
		}
		stream.Observer.OnComplete()
	})

	stream.Publish().Subscribe(func(data []byte) {
		log.Print(data)
	})
  ```
