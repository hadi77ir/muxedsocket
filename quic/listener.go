package quic

import (
	"context"
	"crypto/tls"
	Q "github.com/lucas-clemente/quic-go"
	"muxedsocket"
	"net"
)

var _ muxedsocket.MuxedListener = &Listener{}

type Listener struct {
	listener Q.Listener
}

func (l *Listener) Close() error {
	return l.listener.Close()
}

func (l *Listener) Addr() net.Addr {
	return l.listener.Addr()
}

func (l *Listener) Accept() (muxedsocket.Socket, error) {
	conn, err := l.listener.Accept(context.Background())
	if err != nil {
		return nil, err
	}
	return wrapConn(conn)
}

func Listen(packetConn net.PacketConn, tlsConfig *tls.Config, params *muxedsocket.ServerParams) (muxedsocket.MuxedListener, error) {
	listener, err := Q.Listen(packetConn, tlsConfig, getConfig(params.CommonParams))
	if err != nil {
		return nil, err
	}
	return wrapListener(listener), nil
}
