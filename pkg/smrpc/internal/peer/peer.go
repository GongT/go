package peer

import "net"

type Peer struct {
	Socket *net.TCPConn

	name string

	errCh  chan error
	doneCh chan struct{}
	closed bool
}

func (p *Peer) InitPeer(name string) {
	p.name = name
	p.errCh = make(chan error)
	p.doneCh = make(chan struct{})
}

func (p *Peer) Close() {
	if p.closed {
		return
	}
	p.closed = true
	close(p.doneCh)
	if p.Socket != nil {
		p.Socket.Close()
	}
}

func (p *Peer) Name() string               { return p.name }
func (p *Peer) CloseChan() <-chan struct{} { return p.doneCh }
func (p *Peer) IsClosed() bool             { return p.closed }
func (p *Peer) ErrChan() <-chan error      { return p.errCh }

// TriggerIgnore 把错误发送到 ErrChan 通道中，表示忽略该错误
func (p *Peer) TriggerIgnore(err error) {
	p.errCh <- err
}

func (p *Peer) Call() {
}
