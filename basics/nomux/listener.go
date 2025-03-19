package nomux

import (
	"github.com/hadi77ir/muxedsocket/types"
	"net"
)

type Listener struct {
	listener types.StreamListener
}

func (l *Listener) CloseChan() <-chan struct{} {
	return l.listener.CloseChan()
}

func (l *Listener) Close() error {
	return l.listener.Close()
}

func (l *Listener) Accept() (socket types.Socket, err error) {
	return l.AcceptMuxed()
}

func (l *Listener) Addr() net.Addr {
	return l.listener.Addr()
}

func (l *Listener) AcceptMuxed() (socket types.MuxedSocket, err error) {
	conn, err := l.listener.AcceptConn()
	if err != nil {
		return nil, err
	}
	return WrapAcceptedConn(conn), nil
}

func WrapServer(conn types.StreamListenFunc) (types.MuxedListener, error) {
	listener, err := conn()
	if err != nil {
		return nil, err
	}
	return &Listener{listener: listener}, nil
}

var _ types.MuxedListener = &Listener{}
