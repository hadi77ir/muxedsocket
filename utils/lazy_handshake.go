package utils

import (
	"context"
	"net"
	"sync/atomic"
	"time"
)

type LazyHandshakeConn struct {
	net.Conn
	handshakeFn      func(ctx context.Context) error
	handshakeStarted atomic.Bool
	handshakeDone    chan struct{}
	handshakeContext context.Context
}

func (c *LazyHandshakeConn) SetDeadline(t time.Time) error {
	c.handshakeContext, _ = context.WithDeadline(context.Background(), t)
	return c.Conn.SetDeadline(t)
}

func (c *LazyHandshakeConn) Read(b []byte) (n int, err error) {
	if err := c.guardedHandshake(); err != nil {
		return 0, err
	}
	return c.Conn.Read(b)
}

func (c *LazyHandshakeConn) Write(b []byte) (n int, err error) {
	if err := c.guardedHandshake(); err != nil {
		return 0, err
	}
	return c.Conn.Write(b)
}

func (c *LazyHandshakeConn) guardedHandshake() error {
	oldState := c.handshakeStarted.Swap(true)
	if oldState == false {
		return c.handshakeFn(c.handshakeContext)
	}
	<-c.handshakeDone
	return nil
}

var _ net.Conn = &LazyHandshakeConn{}

func WrapLazyHandshakingConn(conn net.Conn, handshakeFn func(ctx context.Context) error) net.Conn {
	return &LazyHandshakeConn{Conn: conn, handshakeFn: handshakeFn, handshakeDone: make(chan struct{}, 1), handshakeContext: context.Background()}
}
