package packet

import (
	"github.com/hadi77ir/muxedsocket"
	"github.com/hadi77ir/muxedsocket/types"
	"github.com/hadi77ir/muxedsocket/utils"
	"net"
)

type basicWrapper struct {
	net.PacketConn
	remoteAddr net.Addr
	closed     chan struct{}
	dialFunc   StandardPrimedPacketConnFunc
	afterDial  WrappedHookFunc
}

func (c *basicWrapper) CloseChan() <-chan struct{} {
	return c.closed
}

func (c *basicWrapper) Close() error {
	select {
	case <-c.closed:
		break
	default:
		close(c.closed)
	}
	return c.Close()
}

func (c *basicWrapper) RemoteAddr() net.Addr {
	return c.remoteAddr
}

func (c *basicWrapper) CanRedial() bool {
	return c.dialFunc != nil
}

func (c *basicWrapper) Redial() (types.Socket, error) {
	if c.dialFunc == nil {
		return nil, muxedsocket.ErrRedialNotSupported
	}
	conn, err := c.dialFunc()
	if err != nil {
		return nil, err
	}
	return WrapConn(conn, c.remoteAddr, c.dialFunc, c.afterDial), nil
}

func (c *basicWrapper) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	n, addr, err = c.PacketConn.ReadFrom(p)
	c.handleError(err)
	return
}

func (c *basicWrapper) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	n, err = c.PacketConn.WriteTo(p, addr)
	c.handleError(err)
	return
}

func (c *basicWrapper) handleError(err error) {
	if utils.IsConnEOL(err) {
		_ = c.Close()
	}
}

var _ types.PacketConn = &basicWrapper{}

func wrapRawConn(conn net.PacketConn, remoteAddr net.Addr, dialFunc StandardPrimedPacketConnFunc, afterDial WrappedHookFunc) types.PacketConn {
	return &basicWrapper{closed: make(chan struct{}, 1), PacketConn: conn, remoteAddr: remoteAddr, dialFunc: dialFunc, afterDial: afterDial}
}
