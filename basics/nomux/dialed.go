package nomux

import (
	"github.com/hadi77ir/muxedsocket/types"
	"net"
	"sync"
)

func WrapDialedConn(dialFunc types.StreamDialFunc) (types.MuxedSocket, error) {
	return &DialedConn{dialer: dialFunc, closed: make(chan struct{}, 1)}, nil
}

type DialedConn struct {
	dialer          types.StreamDialFunc
	closed          chan struct{}
	remoteAddr      net.Addr
	remoteAddrMutex sync.Mutex
}

func (c *DialedConn) CloseChan() <-chan struct{} {
	return c.closed
}

func (c *DialedConn) Close() error {
	select {
	case <-c.closed:
		return nil
	default:
	}
	close(c.closed)
	return nil
}

func (c *DialedConn) LocalAddr() net.Addr {
	return types.EmptyAddr("nomux:local")
}

func (c *DialedConn) RemoteAddr() net.Addr {
	if c.remoteAddr != nil {
		return c.remoteAddr
	}
	return types.EmptyAddr("nomux:remote")
}

func (c *DialedConn) CanRedial() bool {
	return true
}

func (c *DialedConn) Redial() (types.Socket, error) {
	return WrapDialedConn(c.dialer)
}

func (c *DialedConn) AcceptStream() (stream types.MuxStream, err error) {
	return c.DialStream()
}

func (c *DialedConn) OpenStream() (stream types.MuxStream, err error) {
	return c.DialStream()
}

func (c *DialedConn) DialStream() (stream types.MuxStream, err error) {
	conn, err := c.dialer()
	if err != nil {
		return
	}
	stream = &Stream{StreamConn: conn, dialer: c.DialStream}
	c.remoteAddrMutex.Lock()
	c.remoteAddr = stream.RemoteAddr()
	c.remoteAddrMutex.Unlock()
	// close stream if nomux instance was destroyed
	go func() {
		select {
		case <-c.closed:
		case <-stream.CloseChan():
		}
		_ = stream.Close()
	}()
	return
}
