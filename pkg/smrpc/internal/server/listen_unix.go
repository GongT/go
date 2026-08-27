//go:build win32_socket || !windows

package server

import (
	"errors"
	"net"
	"os"
	"syscall"
)

func (s *Server) listen() (net.Listener, error) {
	socketPath := s.SocketPath()
	addr := &net.UnixAddr{Name: socketPath, Net: "unix"}

	listener, err := net.ListenUnix("unix", addr)
	if err != nil {
		if errors.Is(err, syscall.EADDRINUSE) {
			conn, err := net.DialUnix("unix", nil, addr)
			if err == nil {
				conn.Close()
				return nil, errors.New("其他进程正在使用该RPC管道，无法启动")
			}
			os.Remove(socketPath)
			listener, err = net.ListenUnix("unix", addr)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return listener, nil
}
