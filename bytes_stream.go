package stream

// BytesStream  byte데이터만 흐르는 stream
type BytesStream struct {
	stream      chan []byte
	Observer    *Observer
	AtSubscribe func(data []byte)
}

// NewBytesStream ByteStream 생성
func NewBytesStream(handler *ObservHandler) *BytesStream {
	bs := &BytesStream{
		stream:   make(chan []byte, 1),
		Observer: NewObserver(handler),
	}

	return bs
}

// Send 데이터를 전송한다.
func (bs *BytesStream) Send(data []byte) {
	bs.stream <- data
}

// Publish 관찰중인 stream을 발행한다.
func (bs *BytesStream) Publish() *BytesStream {
	bs.Observer.Observ()

	return bs
}

// Subscribe 데이터를 구독
func (bs *BytesStream) Subscribe(call func([]byte)) {
	if call != nil {
		bs.AtSubscribe = call
	}

SubLoop:
	for {
		select {
		case <-bs.Observer.DoneSubscribe:
			break SubLoop
		case data := <-bs.stream:
			bs.AtSubscribe(data)
		}
	}
}

// Cancel Observer Cancel메소드의 프록시 메소드
func (bs *BytesStream) Cancel() {
	bs.Observer.Cancel()
}
