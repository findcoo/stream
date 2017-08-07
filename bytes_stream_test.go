package stream

import (
	"log"
	"strconv"
	"testing"
)

func TestBytesStream(t *testing.T) {
	handler := &ObservHandler{
		AtCancel:   func() {},
		AtComplete: func() {},
		AtError:    func(error) {},
	}

	obv := NewObserver(handler)
	stream := NewBytesStream(obv)
	stream.Observable = func() {
		for i := 0; i <= 10; i++ {
			stream.Send([]byte{byte(i)})
		}
		stream.OnComplete()
	}

	stream.Publish(nil).Subscribe(func(data []byte) {
		log.Print(data)
	})
}

func TestCancel(t *testing.T) {
	handler := &ObservHandler{
		AtCancel:   func() {},
		AtComplete: func() {},
		AtError:    func(error) {},
	}

	obv := NewObserver(handler)
	stream := NewBytesStream(obv)
	stream.Observable = func() {
		for i := 0; i <= 100; i++ {
			select {
			case <-stream.AfterCancel():
				return
			default:
				stream.Send([]byte(strconv.Itoa(i)))
			}
		}
		stream.OnComplete()
	}

	stream.Cancel()
	stream.Publish(nil).Subscribe(func(data []byte) {
		t.Fail()
	})
}
