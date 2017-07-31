package stream

// BytesStream  byte데이터만 흐르는 stream
type BytesStream struct {
	stream       chan []byte
	Observer     *Observer
	Subscribable func(data []byte)
}

// NewByteStream ByteStream 생성
func NewByteStream() *BytesStream {
	bs := &BytesStream{
		stream:   make(chan []byte, 1),
		Observer: NewObserver(),
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
func (bs *BytesStream) Subscribe(callArray ...func([]byte)) {
SubLoop:
	for {
		select {
		case <-bs.Observer.doneSubscribe:
			break SubLoop
		case data := <-bs.stream:
			if len(callArray) != 0 {
				for _, call := range callArray {
					call(data)
				}
			} else {
				bs.Subscribable(data)
			}
		}
	}
}

// Cancel Observer Cancel메소드의 프록시 메소드
func (bs *BytesStream) Cancel() {
	bs.Observer.Cancel()
}
