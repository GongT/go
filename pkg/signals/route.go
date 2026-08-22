package signals

import (
	"os"
	"os/signal"
	"syscall"
)

type SignalRouter struct {
	signals []os.Signal

	signalChan chan os.Signal
}

func NewRouter(signals ...syscall.Signal) *SignalRouter {
	r := &SignalRouter{
		signalChan: make(chan os.Signal, 1),
	}

	r.signals = make([]os.Signal, len(signals))
	for i, s := range signals {
		r.signals[i] = s
	}
	if len(r.signals) == 0 {
		panic("SignalRouter至少需要一个信号")
	}

	signal.Notify(r.signalChan, r.signals...)

	return r
}

func (r *SignalRouter) Destroy(sig os.Signal) {
	signal.Stop(r.signalChan)
	close(r.signalChan)
}

func (r *SignalRouter) WaitForInt() <-chan struct{} {
	return r.WaitFor(syscall.SIGINT)
}

func (r *SignalRouter) WaitForTerm() <-chan struct{} {
	return r.WaitFor(syscall.SIGTERM)
}

func (r *SignalRouter) WaitFor(signal syscall.Signal) <-chan struct{} {
	// TODO 检查singal是否在r.signals中
	ch := make(chan struct{})
	go func() {
		for {
			s := <-r.signalChan
			if s.(syscall.Signal) == signal {
				close(ch)
				break
			}
		}
	}()
	return ch
}
