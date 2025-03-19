package yamux

import (
	"context"
	"github.com/hadi77ir/muxedsocket"
	"github.com/hadi77ir/muxedsocket/types"
	Y "github.com/hashicorp/yamux"
	"net"
)

var _ types.MuxedSocket = &Conn{}

type Conn struct {
	session *Y.Session
	ctx     context.Context
}

func (c *Conn) CloseChan() <-chan struct{} {
	return c.session.CloseChan()
}

// ClientMuxer dials the target server and establishes a connection. If clientParams.TLSConfig is not nil, a TLS layer is added
// on top of the connection, else not.
func ClientMuxer(conn net.Conn, clientParams *muxedsocket.ClientParams) (types.MuxedSocket, error) {
	sclient, err := Y.Client(conn, getConfig(clientParams.CommonParams))
	if err != nil {
		return nil, err
	}
	return wrapConn(sclient), nil
}

func (c *Conn) AcceptStream() (stream types.MuxStream, err error) {
	s, err := c.session.AcceptStream()
	if err != nil {
		return nil, err
	}
	return wrapStream(s, c.session), nil
}

func (c *Conn) OpenStream() (stream types.MuxStream, err error) {
	s, err := c.session.OpenStream()
	if err != nil {
		return nil, err
	}
	return wrapStream(s, c.session), nil
}

func (c *Conn) Close() error {
	return c.session.Close()
}

func (c *Conn) LocalAddr() net.Addr {
	return c.session.LocalAddr()
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.session.RemoteAddr()
}
