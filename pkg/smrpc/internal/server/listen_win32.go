//go:build !win32_socket && windows

package server

import (
	"net"

	"github.com/Microsoft/go-winio"
)

func (s *Server) listen() (net.Listener, error) {
	pipePath := s.SocketPath()

	config := &winio.PipeConfig{
		InputBufferSize:  4096,
		OutputBufferSize: 4096,
	}

	return winio.ListenPipe(pipePath, config)
}
