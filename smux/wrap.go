package smux

import (
	"github.com/hadi77ir/muxedsocket"
	"github.com/xtaci/smux"
	"net"
)

func wrapConn(conn *smux.Session) muxedsocket.MuxedSocket {
	return &Conn{
		session: conn,
	}
}

func wrapListener(listener net.Listener, config *smux.Config) *Listener {
	return &Listener{listener: listener, config: config}
}

func wrapStream(stream *smux.Stream, session *smux.Session) muxedsocket.MuxStream {
	return &Stream{
		stream:     stream,
		localAddr:  muxedsocket.WrapAddr(session.LocalAddr(), int(stream.ID())),
		remoteAddr: muxedsocket.WrapAddr(session.RemoteAddr(), int(stream.ID())),
	}
}
