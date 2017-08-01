package stream

import (
	"log"
	"strconv"
	"testing"
)

func TestBytesStream(t *testing.T) {
	handler := ObservHandler{
		AtCancel:   func() {},
		AtComplete: func() {},
		AtError:    func(error) {},
	}

	stream := NewByteStream(handler)
	stream.Observer.SetObservable(func() {
		for i := 0; i <= 10; i++ {
			stream.Send([]byte{byte(i)})
		}
		stream.Observer.OnComplete()
	})

	stream.Publish().Subscribe(func(data []byte) {
		log.Print(data)
	})
}

func TestCancel(t *testing.T) {
	handler := ObservHandler{
		AtCancel:   func() {},
		AtComplete: func() {},
		AtError:    func(error) {},
	}

	stream := NewByteStream(handler)
	stream.Observer.SetObservable(func() {
		for i := 0; i <= 100; i++ {
			select {
			case <-stream.Observer.AfterCancel():
				return
			default:
				stream.Send([]byte(strconv.Itoa(i)))
			}
		}
		stream.Observer.OnComplete()
	})

	stream.Cancel()

	stream.Publish().Subscribe(func(data []byte) {
		t.Fail()
	})
}
