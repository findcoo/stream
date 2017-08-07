package stream

// BytesStream  byte데이터만 흐르는 stream
type BytesStream struct {
	stream      chan []byte
	AtSubscribe func(data []byte)
	*Observer
}

// NewBytesStream ByteStream 생성
func NewBytesStream(obv *Observer) *BytesStream {
	bs := &BytesStream{
		stream:   make(chan []byte, 1),
		Observer: obv,
	}

	return bs
}

// Send 데이터를 전송한다.
func (bs *BytesStream) Send(data []byte) {
	bs.WG.Add(1)
	bs.stream <- data
}

// Publish 관찰중인 stream을 발행한다.
func (bs *BytesStream) Publish(target func()) *BytesStream {
	bs.Watch(target)

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
		case <-bs.DoneSubscribe:
			break SubLoop
		case data := <-bs.stream:
			bs.AtSubscribe(data)
			bs.WG.Done()
		}
	}
}
