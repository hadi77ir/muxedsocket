package yamux

import (
	"github.com/hadi77ir/muxedsocket"
	"github.com/hadi77ir/muxedsocket/types"
	Y "github.com/hashicorp/yamux"
	"net"
)

var _ types.MuxedListener = &Listener{}

type Listener struct {
	listener net.Listener
	config   *Y.Config
}

// ServerMuxer takes unencrypted stream-oriented connection-based "listener" and adds TLS and multiplexing.
// If params.TLSConfig is nil, TLS won't be added.
func ServerMuxer(listener net.Listener, params *muxedsocket.ServerParams) (types.MuxedListener, error) {
	return wrapListener(listener, getConfig(params.CommonParams)), nil
}

func (l *Listener) Accept() (socket types.Socket, err error) {
	conn, err := l.listener.Accept()
	if err != nil {
		return nil, err
	}

	yconn, err := Y.Server(conn, l.config)
	if err != nil {
		return nil, err
	}
	return wrapConn(yconn), nil
}

func (l *Listener) Close() error {
	return l.listener.Close()
}

func (l *Listener) Addr() net.Addr {
	return l.listener.Addr()
}
