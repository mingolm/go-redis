package pool

import (
	"bufio"
	"context"
	"errors"
	"io"
	"net"
	"syscall"
	"time"
)

type connTyp uint8

const (
	connTypKeepalive connTyp = iota // 持久化
	connTypTmp                      // 临时
)

type Conn struct {
	netConn   net.Conn
	reader    *bufio.Reader
	writer    *bufio.Writer
	typ       connTyp
	createdAt time.Time
	usedAt    time.Time
}

func NewConnect(conn net.Conn) *Conn {
	return &Conn{
		netConn:   conn,
		reader:    bufio.NewReader(conn),
		writer:    bufio.NewWriter(conn),
		createdAt: time.Now(),
		usedAt:    time.Now(),
	}
}

func (c *Conn) WithWrite(ctx context.Context, wf func(context.Context, *bufio.Writer) error) error {
	defer func() {
		c.writer.Reset(c.netConn)
	}()

	if err := wf(ctx, c.writer); err != nil {
		return err
	}
	if err := c.writer.Flush(); err != nil {
		return err
	}
	return nil
}

func (c *Conn) WithRead(ctx context.Context, rf func(context.Context, *bufio.Reader) error) error {
	defer func() {
		c.reader.Reset(c.netConn)
	}()

	if err := rf(ctx, c.reader); err != nil {
		return err
	}
	return nil
}

func (c *Conn) check() error {
	// Reset previous timeout.
	_ = c.netConn.SetDeadline(time.Time{})

	sysConn, ok := c.netConn.(syscall.Conn)
	if !ok {
		return nil
	}
	rawConn, err := sysConn.SyscallConn()
	if err != nil {
		return err
	}

	var sysErr error

	if err := rawConn.Read(func(fd uintptr) bool {
		var buf [1]byte
		n, err := syscall.Read(int(fd), buf[:])
		switch {
		case n == 0 && err == nil:
			sysErr = io.EOF
		case n > 0:
			sysErr = errors.New("unexpected read from socket")
		case err == syscall.EAGAIN || err == syscall.EWOULDBLOCK:
			sysErr = nil
		default:
			sysErr = err
		}
		return true
	}); err != nil {
		return err
	}

	return sysErr
}
