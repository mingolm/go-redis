package pool

import (
	"bufio"
	"context"
	"net"
	"time"
)

func NewConn(netConn net.Conn) *Conn {
	return &Conn{
		netConn:   netConn,
		reader:    bufio.NewReader(netConn),
		writer:    bufio.NewWriter(netConn),
		createdAt: time.Now(),
	}
}

type Conn struct {
	netConn   net.Conn
	reader    *bufio.Reader
	writer    *bufio.Writer
	createdAt time.Time
}

func (c *Conn) WithWrite(ctx context.Context, wf func(*bufio.Writer) error) error {

	return nil
}

func (c *Conn) WithRead(ctx context.Context, rf func(*bufio.Reader) error) error {
	return nil
}
