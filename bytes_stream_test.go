package stream

import (
	"log"
	"strconv"
	"testing"
)

func TestBytesStream(t *testing.T) {
	stream := NewByteStream()
	stream.Observer.Handler.Observable = func() {
		for i := 0; i <= 10; i++ {
			stream.Send([]byte{byte(i)})
		}
		stream.Observer.OnComplete()
	}

	stream.Publish().Subscribe(func(data []byte) {
		log.Print(data)
	})
}

func TestCancel(t *testing.T) {
	stream := NewByteStream()
	stream.Observer.Handler = Handler{
		Observable: func() {
			for i := 0; i <= 1000; i++ {
				select {
				case <-stream.Observer.AfterCancel():
					return
				default:
					stream.Send([]byte(strconv.Itoa(i)))
				}
			}
			stream.Observer.OnComplete()
		},
		AtCancel:   func() {},
		AtComplete: func() {},
	}

	var i int
	stream.Publish().Subscribe(func(data []byte) {
		log.Print(string(data))
		i++
	})
}
