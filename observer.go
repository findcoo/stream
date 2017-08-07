package stream

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type (
	// StateHandler 상태 변화에 호출되는 함수 형태
	StateHandler func()
	// ErrHandler error 예외처리를 하는 함수 형태
	ErrHandler func(error)
)

// Observer Status를 모니터링하는 관찰자 객체
type Observer struct {
	err           chan error
	cancel        chan struct{}
	DoneSubscribe chan struct{}
	doneObserv    chan struct{}
	Handler       ObservHandler
	Observable    func()
	WG            sync.WaitGroup
}

// ObservHandler Observer의 상태 변화에 따라 실행될 핸들러
type ObservHandler struct {
	AtComplete StateHandler
	AtCancel   StateHandler
	AtError    ErrHandler
}

// DefaultObservHandler 기본 핸들러 설정
func DefaultObservHandler() *ObservHandler {
	return &ObservHandler{
		AtComplete: func() {},
		AtCancel:   func() {},
		AtError:    func(error) {},
	}
}

// AfterSignal process의 종료을 뜻하는 모든 os signal, system call을 감시
func AfterSignal() chan os.Signal {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill, syscall.SIGTERM)
	return sig
}

// NewObserver Observer 생성
func NewObserver(handler *ObservHandler) *Observer {
	if handler == nil {
		handler = DefaultObservHandler()
	}

	obv := &Observer{
		doneObserv:    make(chan struct{}, 1),
		err:           make(chan error, 1),
		cancel:        make(chan struct{}, 1),
		DoneSubscribe: make(chan struct{}, 1),
		Handler:       *handler,
		Observable:    func() {},
	}

	return obv
}

// Watch watch the target handler
func (o *Observer) Watch(target func()) {
	sig := AfterSignal()

	if target != nil {
		o.Observable = target
	}

	go func() {
		go o.Observable()

	ObservLoop:
		for {
			select {
			case state := <-sig:
				log.Printf("capture signal: %s: end observ", state)
				o.endSubscribe()
				return
			case <-o.doneObserv:
				o.Handler.AtComplete()
				o.endSubscribe()
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
func (o *Observer) Cancel() {
	o.cancel <- struct{}{}
	o.Handler.AtCancel()
	o.doneObserv <- struct{}{}
}

// AfterCancel Observable을 취소하기 위한 메소드
func (o *Observer) AfterCancel() <-chan struct{} {
	return o.cancel
}

func (o *Observer) endSubscribe() {
	o.WG.Wait()
	o.DoneSubscribe <- struct{}{}
}
