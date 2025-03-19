package smux

import (
	"github.com/hadi77ir/muxedsocket/types"
	"github.com/xtaci/smux"
	"net"
)

func wrapConn(conn *smux.Session) types.MuxedSocket {
	return &Conn{
		session: conn,
	}
}

func wrapListener(listener net.Listener, config *smux.Config) *Listener {
	return &Listener{listener: listener, config: config}
}

func wrapStream(stream *smux.Stream, session *smux.Session) types.MuxStream {
	return &Stream{
		stream:     stream,
		localAddr:  types.WrapAddr(session.LocalAddr(), int(stream.ID())),
		remoteAddr: types.WrapAddr(session.RemoteAddr(), int(stream.ID())),
	}
}
