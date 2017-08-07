# stream
[![Go Travis Ci](https://travis-ci.org/findcoo/stream.svg?branch=master)](https://travis-ci.org/findcoo/stream)
[![Go Report Card](https://goreportcard.com/badge/github.com/findcoo/stream)](https://goreportcard.com/report/github.com/findcoo/stream)

make simple ReactiveX pattern for golang

# Subsets
* [Kafka](https://github.com/findcoo/stream/tree/master/kafka#kafka)

# Example
* bytes stream
  ```go
	obv := NewObserver(nil)
	stream := NewBytesStream(obv)
	stream.Target = func() {
		for i := 0; i <= 10; i++ {
			stream.Send([]byte{byte(i)})
		}
		stream.OnComplete()
	}

	stream.Publish(nil).Subscribe(func(data []byte) {
		log.Print(data)
	})
  ```
