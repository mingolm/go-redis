package pool

import (
	"bufio"
	"context"
	"net"
	"time"
)

type Conn struct {
	netConn   net.Conn
	reader    *bufio.Reader
	writer    *bufio.Writer
	createdAt time.Time
}

func (c *Conn) WithWrite(ctx context.Context, wf func(*bufio.Writer) error) error {
	if err := wf(c.writer); err != nil {
		return err
	}
	if err := c.writer.Flush(); err != nil {
		return err
	}
	return nil
}

func (c *Conn) WithRead(ctx context.Context, rf func(*bufio.Reader) error) error {
	if err := rf(c.reader); err != nil {
		return err
	}
	c.reader.Reset(c.netConn)
	return nil
}
