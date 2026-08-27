package client

import (
	"log"

	"github.com/gongt/go/pkg/errors"
	"github.com/gongt/go/pkg/smrpc/internal/peer"
)

type Client struct {
	peer.Peer
}

func New(name string) *Client {
	r := &Client{}
	r.InitPeer(name)
	return r
}

func (c *Client) Start() error {
	socket, err := c.connect()
	if err != nil {
		return errors.Extend(err, "无法连接RPC管道，服务是否启动？")
	}

	go func() {
		// 处理接收数据
		var buff = make([]byte, 1024)
		for {
			n, err := socket.Read(buff)
			if err != nil {
				log.Printf("读取数据失败: %v\n", err)
				return
			}
			log.Printf("收到数据: %s\n", string(buff[:n]))
		}
	}()

	return nil
}
