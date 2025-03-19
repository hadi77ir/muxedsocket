package nomux

import (
	"github.com/hadi77ir/muxedsocket"
	"github.com/hadi77ir/muxedsocket/types"
	"net"
	"sync/atomic"
)

type AcceptedConn struct {
	types.StreamConn
	used atomic.Bool
}

func (c *AcceptedConn) CanRedial() bool {
	return false
}

func (c *AcceptedConn) Redial() (types.Socket, error) {
	return nil, muxedsocket.ErrRedialNotSupported
}

func (c *AcceptedConn) AcceptStream() (stream types.MuxStream, err error) {
	return c.UseStream()
}

func (c *AcceptedConn) OpenStream() (stream types.MuxStream, err error) {
	return c.UseStream()
}

func (c *AcceptedConn) UseStream() (stream types.MuxStream, err error) {
	select {
	case <-c.CloseChan():
		return nil, net.ErrClosed
	default:
	}
	oldState := c.used.Swap(true)
	if oldState != true {
		return WrapStream(c.StreamConn, nil), nil
	}
	<-c.StreamConn.CloseChan()
	return nil, net.ErrClosed
}

func WrapAcceptedConn(conn types.StreamConn) types.MuxedSocket {
	return &AcceptedConn{StreamConn: conn}
}

var _ types.MuxedSocket = &AcceptedConn{}
