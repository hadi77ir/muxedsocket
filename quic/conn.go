package quic

import (
	"context"
	"crypto/tls"
	"github.com/hadi77ir/muxedsocket"
	Q "github.com/lucas-clemente/quic-go"
	"net"
)

var _ muxedsocket.DatagramCapableMuxedSocket = &Conn{}

const GracefulCloseCode Q.ApplicationErrorCode = 0x2
const GracefulCloseString string = "graceful close"

type Conn struct {
	conn Q.Connection
}

// Dial dials the target server and establishes a connection.
func Dial(pConn net.PacketConn, tlsConfig *tls.Config, clientParams *muxedsocket.ClientParams) (muxedsocket.MuxedSocket, error) {
	conn, err := Q.Dial(pConn, clientParams.RemoteAddr, tlsConfig.ServerName, tlsConfig, getConfig(clientParams.CommonParams))
	if err != nil {
		return nil, err
	}
	return wrapConn(conn)
}

func (c *Conn) AcceptStream() (stream muxedsocket.MuxStream, err error) {
	s, err := c.conn.AcceptStream(context.Background())
	if err != nil {
		return nil, err
	}
	return wrapStream(c.conn, s)
}

func (c *Conn) OpenStream() (stream muxedsocket.MuxStream, err error) {
	s, err := c.conn.OpenStream()
	if err != nil {
		return nil, err
	}
	return wrapStream(c.conn, s)
}

func (c *Conn) SendDatagram(b []byte) error {
	return c.conn.SendMessage(b)
}

func (c *Conn) ReceiveDatagram() ([]byte, error) {
	return c.conn.ReceiveMessage()
}

func (c *Conn) Close() error {
	return c.conn.CloseWithError(GracefulCloseCode, GracefulCloseString)
}

func (c *Conn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}
