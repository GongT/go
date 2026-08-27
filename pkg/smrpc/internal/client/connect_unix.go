//go:build win32_socket || !windows

package client

import (
	"net"
)

func (c *Client) connect() (net.Conn, error) {
	socketPath := c.SocketPath()
	addr := &net.UnixAddr{Name: socketPath, Net: "unix"}

	return net.DialUnix("unix", nil, addr)
}
