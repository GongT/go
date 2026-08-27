//go:build !win32_socket && windows

package client

import (
	"context"
	"net"

	"github.com/Microsoft/go-winio"
)

func (c *Client) connect() (net.Conn, error) {
	pipePath := c.SocketPath()

	return winio.DialPipeContext(context.Background(), pipePath)
}
