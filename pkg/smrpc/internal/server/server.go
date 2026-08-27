package server

import (
	"log"
	"net"
	"os"

	"github.com/gongt/go/pkg/errors"
	"github.com/gongt/go/pkg/smrpc/internal/peer"
)

type Server struct {
	peer.Peer
}

func New(name string) *Server {
	r := &Server{}
	r.InitPeer(name)
	return r
}

func (s *Server) Start() error {
	listener, err := s.listen()
	if err != nil {
		return errors.Extend(err, "无法监听RPC管道")
	}

	go func() {
		defer os.Remove(s.SocketPath())
		for {
			conn, err := listener.Accept()
			if err != nil {
				if !s.IsClosed() {
					s.TriggerIgnore(err)
				}
				return
			}
			go func() {
				defer conn.Close()
				s.handleConn(conn)
			}()
		}
	}()

	go func() {
		<-s.CloseChan()
		listener.Close()
	}()

	return nil
}

func (s *Server) handleConn(conn net.Conn) {
	log.Println("Accepted connection from", conn.RemoteAddr())
}
