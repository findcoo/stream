package stream

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

type (
	// StateFunc 상태 변화에 호출되는 함수 형태
	StateFunc func()
	// ErrFunc error 예외처리를 하는 함수 형태
	ErrFunc func(error)
)

// Observer Status를 모니터링하는 관찰자 객체
type Observer struct {
	err           chan error
	cancel        chan struct{}
	DoneSubscribe chan struct{}
	doneObserv    chan struct{}
	Handler       Handler
}

// Handler Observer의 상태 변화에 따라 실행될 핸들러
type Handler struct {
	Observable StateFunc
	AtComplete StateFunc
	AtCancel   StateFunc
	AtError    ErrFunc
}

// NewObserver Observer 생성
func NewObserver() *Observer {
	obv := &Observer{
		doneObserv:    make(chan struct{}, 1),
		err:           make(chan error, 1),
		cancel:        make(chan struct{}, 1),
		DoneSubscribe: make(chan struct{}, 1),
		Handler: Handler{
			Observable: func() {},
			AtComplete: func() {},
			AtCancel:   func() {},
			AtError:    func(err error) {},
		},
	}

	return obv
}

// Observ 대상 핸들러를 관찰한다.
func (o *Observer) Observ() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill, syscall.SIGTERM)

	go func() {
		go o.Handler.Observable()

	ObservLoop:
		for {
			select {
			case state := <-sig:
				log.Printf("capture signal: %s: end observ", state)
				o.DoneSubscribe <- struct{}{}
				return
			case <-o.doneObserv:
				o.Handler.AtComplete()
				o.DoneSubscribe <- struct{}{}
				break ObservLoop
			case err := <-o.err:
				o.Handler.AtError(err)
			}
		}
	}()
}

// OnComplete Observer의 활동을 끝낸다.
// Observable 에서 return과 동반되거나 종료될 때 사용
func (o *Observer) OnComplete() {
	o.doneObserv <- struct{}{}
}

// OnError error를 전달한다.
func (o *Observer) OnError(err error) {
	o.err <- err
}

// Cancel Observ 활동을 취소시킨다.
// NOTE Observ의 종료와 동시에 Subscribe도 종료된다.
// NOTE Cancel시에 데이터를 완전히 보장하지 않고, 바로 종료시킨다.
// Observable 에서 return과 동반
func (o *Observer) Cancel() {
	o.cancel <- struct{}{}
	o.Handler.AtCancel()
	o.doneObserv <- struct{}{}
}

// AfterCancel Observable을 취소하기 위한 메소드
func (o *Observer) AfterCancel() <-chan struct{} {
	return o.cancel
}
