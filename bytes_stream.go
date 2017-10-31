package stream

// BytesStream handle data with the sub stream
type BytesPipeline func(data [][]byte) <-chan []byte

// BytesStream is bytestyped go channel stream
type BytesStream struct {
	stream      chan []byte
	AtSubscribe func(data []byte)
	pipelines   []BytesPipeline
	actions     []string
	*Observer
}

// NewBytesStream create the BytesStream
func NewBytesStream(obv *Observer) *BytesStream {
	bs := &BytesStream{
		stream:   make(chan []byte, 1),
		Observer: obv,
	}
	return bs
}

// Send bytes data to the go channel
func (bs *BytesStream) Send(data []byte) {
	if len(bs.pipelines) > 0 {
		mutatedArray := [][]byte{data}
		tempArray := [][]byte{}

		for _, pipe := range bs.pipelines {
			for res := range pipe(mutatedArray) {
				tempArray = append(tempArray, res)
			}
			mutatedArray = tempArray
			tempArray = tempArray[:0]
		}

		for _, item := range mutatedArray {
			bs.WG.Add(1)
			bs.stream <- item
		}
	} else {
		bs.WG.Add(1)
		bs.stream <- data
	}
}

// Publish the stream with the observer
func (bs *BytesStream) Publish(target func()) *BytesStream {
	bs.Watch(target)
	return bs
}

// Chunk gather data until fill the buffer
func (bs *BytesStream) Chunk(size int) *BytesStream {
	pipe := func(payload [][]byte) <-chan []byte {
		line := make(chan []byte)
		go func() {
			for _, data := range payload {
				dataLength := len(data)
				chunkSize := int(dataLength / size)
				surplusSize := dataLength % size

				for i := 0; i < chunkSize; i++ {
					line <- data[i : i+size]
				}
				if surplusSize > 0 {
					line <- data[chunkSize*size:]
				}
			}
			close(line)
		}()
		return line
	}
	bs.pipelines = append(bs.pipelines, pipe)
	return bs
}

// Subscribe data from the go channel
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
