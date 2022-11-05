package smux

import (
	"github.com/hadi77ir/muxedsocket"
	S "github.com/xtaci/smux"
	"net"
)

var _ muxedsocket.MuxedListener = &Listener{}

type Listener struct {
	listener net.Listener
	config   *S.Config
}

// ServerMuxer takes unencrypted stream-oriented connection-based "listener" and adds TLS and multiplexing.
// If params.TLSConfig is nil, TLS won't be added.
func ServerMuxer(listener net.Listener, params *muxedsocket.ServerParams) (muxedsocket.MuxedListener, error) {
	return wrapListener(listener, getConfig(params.CommonParams)), nil
}

func (l *Listener) Accept() (socket muxedsocket.Socket, err error) {
	conn, err := l.listener.Accept()
	if err != nil {
		return nil, err
	}

	sconn, err := S.Server(conn, l.config)
	if err != nil {
		return nil, err
	}
	return wrapConn(sconn), nil
}

func (l *Listener) Close() error {
	return l.listener.Close()
}

func (l *Listener) Addr() net.Addr {
	return l.listener.Addr()
}
