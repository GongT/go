package signals

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type SignalHandler struct {
	signals    []os.Signal
	signalChan chan os.Signal

	stopChan chan uint
	killChan chan bool

	counter uint
}

func NewHandler(signal ...syscall.Signal) *SignalHandler {
	h := &SignalHandler{}

	if len(signal) == 0 {
		signal = []syscall.Signal{syscall.SIGINT, syscall.SIGTERM}
	}

	h.signals = make([]os.Signal, len(signal))
	for i, s := range signal {
		h.signals[i] = s
	}

	return h
}

func (h *SignalHandler) Start() {
	h.signalChan = make(chan os.Signal, 5)
	h.killChan = make(chan bool)
	h.stopChan = make(chan uint)

	signal.Notify(h.signalChan, h.signals...)
}

func (h *SignalHandler) Stop() {
	signal.Stop(h.signalChan)
	close(h.signalChan)
	close(h.stopChan)
	close(h.killChan)
}

func (h *SignalHandler) Wait() {
	for {
		s := <-h.signalChan
		h.counter++

		c := h.counter

		if h.counter >= 3 {
			log.Println("收到3次中断信号，强制退出")
			select {
			case <-time.After(5 * time.Second):
				log.Println("程序未能在5秒内退出")
			case h.killChan <- true:
				// 成功
			}
			os.Exit(1)
		} else if c == 1 {
			log.Printf("收到中断信号%s，准备退出……", s)
		} else {
			log.Printf("收到%d次中断信号%s", c, s)
		}

		select {
		case h.stopChan <- c - 1:
		default:
			log.Println("无人接收中断信号？")
		}
	}
}

// 如果连续收到3次中断信号，则此通道会输出一个true并关闭
// 如果调用了Stop()，则此通道会直接关闭
func (h *SignalHandler) Killed() <-chan bool {
	return h.killChan
}

// 如果收到中断信号，则此通道会输出一个uint，表示收到的中断信号次数
// 如果调用了Stop()，则此通道会直接关闭
func (h *SignalHandler) Done() <-chan uint {
	return h.stopChan
}
