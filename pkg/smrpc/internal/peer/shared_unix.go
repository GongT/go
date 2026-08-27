//go:build !windows

package peer

import (
	"fmt"
)

func (p *Peer) SocketPath() string {
	return fmt.Sprintf(`/tmp/%s.sock`, p.name)
}
